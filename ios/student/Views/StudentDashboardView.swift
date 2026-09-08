import SwiftUI


public struct StudentDashboardView: View {
    @ObservedObject public var viewModel: StudentViewModel
    public var token: String
    public var onLogout: (() -> Void)?
    @State private var studentName: String
    @State private var selectedTab = 0

    public init(token: String = "", studentName: String = "Studente", viewModel: StudentViewModel = StudentViewModel(), onLogout: (() -> Void)? = nil) {
        self.token = token
        self.onLogout = onLogout
        self._studentName = State(initialValue: studentName)
        self.viewModel = viewModel
    }

    public var body: some View {
        TabView(selection: $selectedTab) {
            StudentHomeView(studentName: studentName, viewModel: viewModel, onLogout: onLogout)
                .tabItem {
                    Label(NSLocalizedString("tab_home", comment: ""), systemImage: "house.fill")
                }
                .tag(0)

            StudentGradesView(viewModel: viewModel)
                .tabItem {
                    Label(NSLocalizedString("tab_grades", comment: ""), systemImage: "chart.bar.doc.horizontal.fill")
                }
                .tag(1)

            StudentAgendaView(viewModel: viewModel)
                .tabItem {
                    Label(NSLocalizedString("tab_agenda", comment: ""), systemImage: "calendar")
                }
                .tag(2)

            StudentAttendanceView(viewModel: viewModel)
                .tabItem {
                    Label(NSLocalizedString("tab_attendance", comment: ""), systemImage: "checkmark.circle.fill")
                }
                .tag(3)

            StudentReportCardView(viewModel: viewModel)
                .tabItem {
                    Label(NSLocalizedString("tab_report_card", comment: ""), systemImage: "doc.text.fill")
                }
                .tag(4)
        }
        .tint(Color.purple)
        .task {
            if !token.isEmpty {
                await viewModel.loadFromDatabase(token: token)
            }
        }
    }
}

// 1. Home View
struct StudentHomeView: View {
    let studentName: String
    @ObservedObject var viewModel: StudentViewModel
    var onLogout: (() -> Void)? = nil

    var body: some View {
        NavigationView {
            ScrollView {
                VStack(spacing: 20) {
                    VStack(alignment: .leading, spacing: 12) {
                        HStack {
                            VStack(alignment: .leading, spacing: 4) {
                                Text(String(format: NSLocalizedString("welcome_student", comment: ""), studentName))
                                    .font(.title)
                                    .fontWeight(.bold)
                                    .foregroundColor(.white)
                                Text("\(NSLocalizedString("gpa_average", comment: "")): \(String(format: "%.1f", viewModel.calculateGPA()))")
                                    .font(.subheadline)
                                    .padding(.horizontal, 10)
                                    .padding(.vertical, 4)
                                    .background(Color.white.opacity(0.25))
                                    .cornerRadius(8)
                                    .foregroundColor(.white)
                            }
                            Spacer()
                            Image(systemName: "sparkles")
                                .font(.system(size: 32))
                                .foregroundColor(.white.opacity(0.9))
                        }
                    }
                    .padding(20)
                    .background(StudentTheme.primaryGradient)
                    .cornerRadius(20)
                    .shadow(color: Color.indigo.opacity(0.3), radius: 10, x: 0, y: 5)

                    VStack(alignment: .leading, spacing: 12) {
                        Text(NSLocalizedString("grades_title", comment: ""))
                            .font(.headline)
                            .fontWeight(.bold)
                            .accessibilityAddTraits(.isHeader)

                        ForEach(viewModel.grades) { grade in
                            HStack {
                                VStack(alignment: .leading, spacing: 4) {
                                    Text(grade.subject)
                                        .font(.body)
                                        .fontWeight(.semibold)
                                    Text("\(grade.localizedType) • \(grade.date)")
                                        .font(.caption)
                                        .foregroundColor(.secondary)
                                }
                                Spacer()
                                Text(String(format: "%.1f", grade.grade))
                                    .font(.headline)
                                    .fontWeight(.bold)
                                    .foregroundColor(.white)
                                    .padding(.horizontal, 12)
                                    .padding(.vertical, 6)
                                    .background(grade.grade >= 6.0 ? Color.green : Color.red)
                                    .cornerRadius(10)
                            }
                            .padding()
                            .background(Color(UIColor.secondarySystemGroupedBackground))
                            .cornerRadius(14)
                        }
                    }
                }
                .padding()
            }
            .navigationTitle(NSLocalizedString("dashboard_title", comment: ""))
            .toolbar {
                ToolbarItem(placement: .primaryAction) {
                    if let onLogout = onLogout {
                        Button(action: onLogout) {
                            HStack(spacing: 4) {
                                Image(systemName: "rectangle.portrait.and.arrow.right")
                                Text(NSLocalizedString("logout", comment: ""))
                                    .font(.caption)
                                    .fontWeight(.semibold)
                            }
                            .foregroundColor(.white)
                            .padding(.horizontal, 10)
                            .padding(.vertical, 6)
                            .background(Color.red.opacity(0.85))
                            .cornerRadius(8)
                        }
                    }
                }
            }
            .background(Color(UIColor.systemGroupedBackground).ignoresSafeArea())
        }
    }
}

// 2. Grades View
struct StudentGradesView: View {
    @ObservedObject var viewModel: StudentViewModel

    var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("grades_title", comment: ""))) {
                    HStack {
                        Text(NSLocalizedString("gpa_average", comment: ""))
                            .fontWeight(.semibold)
                        Spacer()
                        Text(String(format: "%.1f", viewModel.calculateGPA()))
                            .fontWeight(.bold)
                            .foregroundColor(.purple)
                    }
                }

                Section(header: Text(NSLocalizedString("grades_title", comment: ""))) {
                    if viewModel.grades.isEmpty {
                        Text("Nessun voto registrato")
                            .font(.subheadline)
                            .foregroundColor(.secondary)
                            .padding(.vertical, 8)
                    } else {
                        ForEach(viewModel.grades) { grade in
                            HStack {
                                VStack(alignment: .leading) {
                                    Text(grade.subject)
                                        .fontWeight(.semibold)
                                    Text("\(grade.localizedType) • \(grade.date)")
                                        .font(.caption)
                                        .foregroundColor(.secondary)
                                }
                                Spacer()
                                Text(String(format: "%.1f", grade.grade))
                                    .fontWeight(.bold)
                                    .foregroundColor(grade.grade >= 6 ? .green : .red)
                            }
                        }
                    }
                }
            }
            .navigationTitle(NSLocalizedString("grades_title", comment: ""))
        }
    }
}

// 3. Agenda View
struct StudentAgendaView: View {
    @ObservedObject var viewModel: StudentViewModel

    var body: some View {
        NavigationView {
            List {
                if viewModel.homework.isEmpty {
                    Text("Nessun compito assegnato in agenda")
                        .font(.subheadline)
                        .foregroundColor(.secondary)
                        .padding(.vertical, 8)
                } else {
                    ForEach(viewModel.homework) { task in
                        HStack {
                            Image(systemName: task.isCompleted ? "checkmark.circle.fill" : "circle")
                                .foregroundColor(task.isCompleted ? .green : .gray)
                                .onTapGesture {
                                    _ = viewModel.toggleHomework(id: task.id)
                                }
                            VStack(alignment: .leading) {
                                Text(task.subject)
                                    .fontWeight(.semibold)
                                Text(task.taskDescription)
                                    .font(.caption)
                                    .foregroundColor(.secondary)
                                Text("\(NSLocalizedString("agenda_title", comment: "")): \(task.dueDate)")
                                    .font(.caption2)
                                    .foregroundColor(.purple)
                            }
                        }
                        .padding(.vertical, 4)
                    }
                }
            }
            .navigationTitle(NSLocalizedString("agenda_title", comment: ""))
        }
    }
}

// 4. Attendance View
struct StudentAttendanceView: View {
    @ObservedObject var viewModel: StudentViewModel

    var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("attendance_title", comment: ""))) {
                    if viewModel.attendance.isEmpty {
                        Text("Nessuna assenza o ritardo registrato")
                            .font(.subheadline)
                            .foregroundColor(.secondary)
                            .padding(.vertical, 8)
                    } else {
                        ForEach(viewModel.attendance) { item in
                            HStack {
                                Text("\(item.date) • \(item.localizedType)")
                                Spacer()
                                if item.isJustified {
                                    Label("Giustificata", systemImage: "checkmark.seal.fill")
                                        .foregroundColor(.green)
                                        .font(.caption)
                                }
                            }
                        }
                    }
                }
            }
            .navigationTitle(NSLocalizedString("attendance_title", comment: ""))
        }
    }
}

// 5. Report Card View
struct StudentReportCardView: View {
    @ObservedObject var viewModel: StudentViewModel

    var body: some View {
        NavigationView {
            ScrollView {
                VStack(spacing: 16) {
                    VStack(alignment: .leading, spacing: 8) {
                        Text(NSLocalizedString("report_card_title", comment: ""))
                            .font(.headline)
                            .fontWeight(.bold)
                            .foregroundColor(.white)
                        Text("\(NSLocalizedString("gpa_average", comment: "")): \(String(format: "%.1f", viewModel.calculateGPA()))")
                            .font(.caption)
                            .foregroundColor(.white.opacity(0.8))
                    }
                    .frame(maxWidth: .infinity, alignment: .leading)
                    .padding()
                    .background(Color.purple)
                    .cornerRadius(16)

                    if viewModel.grades.isEmpty {
                        Text("Nessuna valutazione registrata al momento")
                            .font(.subheadline)
                            .foregroundColor(.secondary)
                            .padding()
                    } else {
                        ForEach(viewModel.grades) { grade in
                            HStack {
                                Text(grade.subject).fontWeight(.medium)
                                Spacer()
                                Text(String(format: "%.1f", grade.grade))
                                    .fontWeight(.bold)
                                    .foregroundColor(.white)
                                    .padding(.horizontal, 12)
                                    .padding(.vertical, 4)
                                    .background(grade.grade >= 6 ? Color.green : Color.red)
                                    .cornerRadius(8)
                            }
                            .padding()
                            .background(Color(UIColor.secondarySystemGroupedBackground))
                            .cornerRadius(12)
                        }
                    }
                }
                .padding()
            }
            .navigationTitle(NSLocalizedString("report_card_title", comment: ""))
            .background(Color(UIColor.systemGroupedBackground).ignoresSafeArea())
        }
    }
}
