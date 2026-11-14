import XCTest
@testable import TeacherApp

private final class StubTeacherAPIService: TeacherAPIServiceProtocol {
    private let classes: [TeacherClassSessionModel]
    private let students: [RollCallStudent]

    init(
        classes: [TeacherClassSessionModel] = [
            TeacherClassSessionModel(id: "class-1", className: "2A", subject: "Matematica")
        ],
        students: [RollCallStudent] = [
            RollCallStudent(id: "student-1", name: "Studente Test")
        ]
    ) {
        self.classes = classes
        self.students = students
    }

    func login(email: String, password: String) async throws -> String {
        "test-token"
    }

    func fetchClasses(token: String) async throws -> [TeacherClassSessionModel] {
        classes
    }

    func fetchClassStudents(token: String, classId: String) async throws -> [RollCallStudent] {
        students
    }

    func signLesson(token: String, classId: String, topic: String) async throws -> TeacherClassSessionModel {
        TeacherClassSessionModel(id: classId, className: "2A", subject: "Matematica", isSigned: true, lessonTopic: topic)
    }

    func submitRollCall(token: String, classId: String, records: [TeacherRollCallRecord]) async throws -> Bool {
        true
    }
}

final class TeacherIntegrationTests: XCTestCase {

    var viewModel: TeacherViewModel!
    var cacheManager: OfflineCacheManager!
    var securityHelper: SecurityCheckHelper!
    var biometricManager: BiometricAuthManager!

    override func setUp() {
        super.setUp()
        viewModel = TeacherViewModel(apiService: StubTeacherAPIService())
        cacheManager = OfflineCacheManager()
        securityHelper = SecurityCheckHelper.shared
        biometricManager = BiometricAuthManager.shared
    }

    func testSecurityAndBiometricIntegrity() {
        let report = securityHelper.getSecurityReport()
        XCTAssertNotNil(report["isJailbroken"])
        XCTAssertNotNil(report["isSimulator"])

        // The result is platform-dependent; evaluating the capability must not crash.
        _ = biometricManager.canEvaluatePolicy()
    }

    func testViewModelLoadsClassesAndStudentsAndCachesThem() async {
        let loaded = await viewModel.loadFromDatabase(token: "test-token")

        XCTAssertTrue(loaded)
        XCTAssertEqual(viewModel.currentSession?.id, "class-1")
        XCTAssertEqual(viewModel.students.map(\.id), ["student-1"])

        guard let currentSession = viewModel.currentSession else {
            XCTFail("The stub service should provide a class")
            return
        }
        cacheManager.saveClasses([currentSession])
        cacheManager.saveStudents(viewModel.students)

        XCTAssertTrue(cacheManager.isCacheValid())
        XCTAssertEqual(cacheManager.getCachedClasses(), [currentSession])
        XCTAssertEqual(cacheManager.getCachedStudents(), viewModel.students)

        cacheManager.clear()
        XCTAssertFalse(cacheManager.isCacheValid())
    }

    func testViewModelHandlesAnAccountWithoutClasses() async {
        viewModel = TeacherViewModel(apiService: StubTeacherAPIService(classes: [], students: []))

        let loaded = await viewModel.loadFromDatabase(token: "test-token")

        XCTAssertTrue(loaded)
        XCTAssertNil(viewModel.currentSession)
        XCTAssertTrue(viewModel.students.isEmpty)
    }
}
