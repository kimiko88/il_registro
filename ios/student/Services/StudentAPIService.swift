import Foundation

public protocol StudentAPIServiceProtocol {
    var lastStudentName: String? { get }
    func login(email: String, password: String) async throws -> String
    func fetchGrades(token: String) async throws -> [GradeItemModel]
    func fetchHomework(token: String) async throws -> [HomeworkItemModel]
    func fetchAttendance(token: String) async throws -> [StudentAttendanceRecordModel]
}

public class HttpStudentAPIService: StudentAPIServiceProtocol {
    private let baseURL: URL
    public var lastStudentName: String?

    public init(baseURL: URL? = nil) {
        if let baseURL = baseURL {
            self.baseURL = baseURL
        } else if let url = URL(string: AppConfig.baseURL) {
            self.baseURL = url
        } else {
            self.baseURL = URL(string: "https://registro-backend-fdu2.onrender.com/api/v1")!
        }
    }

    public func login(email: String, password: String) async throws -> String {
        var cleanEmail = email.trimmingCharacters(in: .whitespacesAndNewlines).lowercased()
        // Trasparenza per alias documentati (es. studentea_1 -> studente2a_1)
        if cleanEmail.hasPrefix("studentea_") {
            cleanEmail = cleanEmail.replacingOccurrences(of: "studentea_", with: "studente2a_")
        } else if cleanEmail.hasPrefix("studenteb_") {
            cleanEmail = cleanEmail.replacingOccurrences(of: "studenteb_", with: "studente2b_")
        }

        let url = baseURL.appendingPathComponent("auth/login")
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")

        let body = ["email": cleanEmail, "password": password]
        request.httpBody = try JSONSerialization.data(withJSONObject: body)

        let (data, response) = try await URLSession.shared.data(for: request)
        guard let httpResponse = response as? HTTPURLResponse, (200...299).contains(httpResponse.statusCode) else {
            if let json = try? JSONSerialization.jsonObject(with: data) as? [String: Any],
               let errMsg = json["error"] as? String ?? json["message"] as? String {
                let code = (response as? HTTPURLResponse)?.statusCode ?? 401
                throw NSError(domain: "StudentAPI", code: code, userInfo: [NSLocalizedDescriptionKey: "\(errMsg) (HTTP \(code))"])
            }
            let code = (response as? HTTPURLResponse)?.statusCode ?? 401
            throw NSError(domain: "StudentAPI", code: code, userInfo: [NSLocalizedDescriptionKey: "Credenziali non valide o errore server (HTTP \(code))"])
        }

        if let json = try JSONSerialization.jsonObject(with: data) as? [String: Any] {
            let userObj = json["user"] as? [String: Any]
            let role = userObj?["role"] as? String
            guard role == nil || role == "student" else {
                throw NSError(domain: "StudentAPI", code: 403, userInfo: [NSLocalizedDescriptionKey: "Accesso non consentito: questo account non appartiene a uno studente."])
            }
            if let firstName = userObj?["first_name"] as? String, let lastName = userObj?["last_name"] as? String {
                self.lastStudentName = "\(firstName) \(lastName)".trimmingCharacters(in: .whitespacesAndNewlines)
            } else if let name = userObj?["first_name"] as? String ?? userObj?["name"] as? String {
                self.lastStudentName = name
            }

            if let token = json["token"] as? String ?? json["access_token"] as? String {
                return token
            }
        }
        throw NSError(domain: "StudentAPI", code: 500, userInfo: [NSLocalizedDescriptionKey: "Token non presente nella risposta"])
    }

    private func fetchSubjectsMap(token: String) async -> [String: String] {
        let url = baseURL.appendingPathComponent("subjects")
        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")

        guard let (data, response) = try? await URLSession.shared.data(for: request),
              let httpResponse = response as? HTTPURLResponse, (200...299).contains(httpResponse.statusCode),
              let jsonArray = try? JSONSerialization.jsonObject(with: data) as? [[String: Any]] else {
            return [:]
        }

        var map: [String: String] = [:]
        for obj in jsonArray {
            if let id = obj["id"] as? String, let name = obj["name"] as? String {
                map[id] = name
            }
        }
        return map
    }

    public func fetchGrades(token: String) async throws -> [GradeItemModel] {
        let subjectMap = await fetchSubjectsMap(token: token)

        let url = baseURL.appendingPathComponent("grades/my-grades")
        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")

        let (data, response) = try await URLSession.shared.data(for: request)
        guard let httpResponse = response as? HTTPURLResponse, (200...299).contains(httpResponse.statusCode) else {
            throw NSError(domain: "StudentAPI", code: 401, userInfo: [NSLocalizedDescriptionKey: "Non autorizzato"])
        }

        var rawGradesList: [[String: Any]] = []

        if let jsonObject = try? JSONSerialization.jsonObject(with: data) {
            if let directArray = jsonObject as? [[String: Any]] {
                rawGradesList = directArray
            } else if let dict = jsonObject as? [String: Any] {
                if let semesters = dict["semesters"] as? [[String: Any]] {
                    for semester in semesters {
                        if let sGrades = semester["grades"] as? [[String: Any]] {
                            rawGradesList.append(contentsOf: sGrades)
                        }
                    }
                } else if let grades = dict["grades"] as? [[String: Any]] {
                    rawGradesList = grades
                }
            }
        }

        return rawGradesList.enumerated().map { index, dict in
            let rawSubjectId = dict["subject_id"] as? String ?? ""
            let subjectName = (dict["subject_name"] as? String) ?? (dict["subject"] as? String) ?? subjectMap[rawSubjectId] ?? "Materia"

            let gradeVal: Double
            if let d = dict["grade_value"] as? Double {
                gradeVal = d
            } else if let i = dict["grade_value"] as? Int {
                gradeVal = Double(i)
            } else if let d = dict["grade"] as? Double {
                gradeVal = d
            } else if let i = dict["grade"] as? Int {
                gradeVal = Double(i)
            } else {
                gradeVal = 0.0
            }

            let weightVal: Double
            if let w = dict["weight"] as? Double {
                weightVal = w
            } else if let i = dict["weight"] as? Int {
                weightVal = Double(i)
            } else {
                weightVal = 1.0
            }

            let rawType = (dict["evaluation_type"] as? String) ?? (dict["grade_type"] as? String) ?? (dict["type"] as? String) ?? "Orale"
            var rawDate = (dict["date"] as? String) ?? ""
            if rawDate.count >= 10 {
                rawDate = String(rawDate.prefix(10))
            }

            return GradeItemModel(
                id: "\(dict["id"] ?? index)",
                subject: subjectName,
                grade: gradeVal,
                weight: weightVal,
                type: rawType,
                date: rawDate
            )
        }
    }

    public func fetchHomework(token: String) async throws -> [HomeworkItemModel] {
        let url = baseURL.appendingPathComponent("agenda")
        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")

        let (data, response) = try await URLSession.shared.data(for: request)
        guard let httpResponse = response as? HTTPURLResponse, (200...299).contains(httpResponse.statusCode) else {
            return []
        }

        guard let jsonArray = try? JSONSerialization.jsonObject(with: data) as? [[String: Any]] else {
            return []
        }

        return jsonArray.enumerated().map { index, dict in
            var dueDate = (dict["due_date"] as? String) ?? (dict["date"] as? String) ?? ""
            if dueDate.count >= 10 {
                dueDate = String(dueDate.prefix(10))
            }
            return HomeworkItemModel(
                id: "\(dict["id"] ?? index)",
                subject: (dict["subject_name"] as? String) ?? (dict["subject"] as? String) ?? (dict["title"] as? String) ?? "Compito",
                taskDescription: (dict["description"] as? String) ?? (dict["notes"] as? String) ?? "",
                dueDate: dueDate,
                isCompleted: (dict["is_completed"] as? Bool) ?? (dict["completed"] as? Bool) ?? false
            )
        }
    }

    public func fetchAttendance(token: String) async throws -> [StudentAttendanceRecordModel] {
        let url = baseURL.appendingPathComponent("attendance/my-attendance")
        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")

        let (data, response) = try await URLSession.shared.data(for: request)
        guard let httpResponse = response as? HTTPURLResponse, (200...299).contains(httpResponse.statusCode) else {
            return []
        }

        guard let jsonArray = try? JSONSerialization.jsonObject(with: data) as? [[String: Any]] else {
            return []
        }

        return jsonArray.enumerated().map { index, dict in
            var rawDate = (dict["date"] as? String) ?? ""
            if rawDate.count >= 10 {
                rawDate = String(rawDate.prefix(10))
            }
            return StudentAttendanceRecordModel(
                id: "\(dict["id"] ?? index)",
                date: rawDate,
                type: (dict["type"] as? String) ?? (dict["status"] as? String) ?? "Assenza",
                isJustified: (dict["is_justified"] as? Bool) ?? (dict["justified"] as? Bool) ?? false
            )
        }
    }
}
