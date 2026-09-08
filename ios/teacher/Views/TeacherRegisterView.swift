import SwiftUI

public struct TeacherRegisterView: View {
    @ObservedObject public var viewModel: TeacherViewModel
    public var token: String
    @State private var selectedTab = 0

    public init(token: String = "", viewModel: TeacherViewModel = TeacherViewModel()) {
        self.token = token
        self.viewModel = viewModel
    }

    public var body: some View {
        TabView(selection: $selectedTab) {
            TeacherFirmaTab(viewModel: viewModel)
                .tabItem { Label(NSLocalizedString("teacher_dashboard_title", comment: ""), systemImage: "pencil.and.list.clipboard") }
                .tag(0)

            TeacherAppelloTab(viewModel: viewModel)
                .tabItem { Label(NSLocalizedString("take_attendance", comment: ""), systemImage: "checkmark.rectangle.stack.fill") }
                .tag(1)

            TeacherGradesTab(viewModel: viewModel)
                .tabItem { Label(NSLocalizedString("add_grade", comment: ""), systemImage: "chart.bar.fill") }
                .tag(2)

            TeacherAgendaTab()
                .tabItem { Label(NSLocalizedString("teacher_dashboard_title", comment: ""), systemImage: "calendar.badge.plus") }
                .tag(3)

            TeacherScrutinyTab()
                .tabItem { Label(NSLocalizedString("scrutiny_board", comment: ""), systemImage: "graduationcap.fill") }
                .tag(4)
        }
        .tint(Color.blue)
        .task {
            if !token.isEmpty {
                _ = await viewModel.loadFromDatabase(token: token)
            }
        }
    }
}

struct TeacherFirmaTab: View {
    @ObservedObject var viewModel: TeacherViewModel

    var body: some View {
        NavigationView {
            ScrollView {
                VStack(spacing: 20) {
                    VStack(alignment: .leading, spacing: 12) {
                        Text(NSLocalizedString("sign_hour", comment: ""))
                            .font(.headline)
                            .fontWeight(.bold)
                        Text("\(NSLocalizedString("select_class", comment: "")) \(viewModel.currentSession?.className ?? "") • \(viewModel.currentSession?.subject ?? "")")
                            .font(.subheadline)
                            .foregroundColor(.secondary)

                        TextField(NSLocalizedString("sign_hour", comment: ""), text: $viewModel.lessonTopic)
                            .textFieldStyle(.roundedBorder)

                        Button(action: { _ = viewModel.signLesson(topic: viewModel.lessonTopic) }) {
                            HStack {
                                Image(systemName: viewModel.isHourSigned ? "checkmark.circle.fill" : "pencil.line")
                                Text(viewModel.isHourSigned ? NSLocalizedString("lesson_signed", comment: "") : NSLocalizedString("sign_hour", comment: ""))
                                    .fontWeight(.semibold)
                            }
                            .frame(maxWidth: .infinity)
                            .padding()
                            .background(viewModel.isHourSigned ? Color.green : Color.blue)
                            .foregroundColor(.white)
                            .cornerRadius(12)
                        }
                    }
                    .padding()
                    .background(Color(UIColor.secondarySystemGroupedBackground))
                    .cornerRadius(16)
                }
                .padding()
            }
            .navigationTitle(NSLocalizedString("teacher_dashboard_title", comment: ""))
        }
    }
}

struct TeacherAppelloTab: View {
    @ObservedObject var viewModel: TeacherViewModel

    var body: some View {
        NavigationView {
            List(viewModel.students) { student in
                HStack {
                    Text(student.name).fontWeight(.medium)
                    Spacer()
                    Button(student.status) {
                        _ = viewModel.toggleAttendance(studentId: student.id, status: student.status == "presente" ? "assente" : "presente")
                    }
                    .buttonStyle(.borderedProminent)
                    .tint(student.status == "presente" ? .green : .red)
                }
            }
            .navigationTitle(NSLocalizedString("take_attendance", comment: ""))
        }
    }
}

struct TeacherGradesTab: View {
    @ObservedObject var viewModel: TeacherViewModel

    var body: some View {
        NavigationView {
            List(viewModel.students) { student in
                HStack {
                    Text(student.name)
                    Spacer()
                    Button(NSLocalizedString("add_grade", comment: "")) {
                        _ = viewModel.insertGrade(studentId: student.id, grade: 8.0, type: "Scritto")
                    }
                    .buttonStyle(.bordered)
                }
            }
            .navigationTitle(NSLocalizedString("add_grade", comment: ""))
        }
    }
}

struct TeacherAgendaTab: View {
    var body: some View {
        NavigationView {
            List {
                Button(action: {}) {
                    Label(NSLocalizedString("teacher_dashboard_title", comment: ""), systemImage: "plus.circle.fill")
                }
            }
            .navigationTitle(NSLocalizedString("teacher_dashboard_title", comment: ""))
        }
    }
}

struct TeacherScrutinyTab: View {
    var body: some View {
        NavigationView {
            ScrollView {
                VStack(spacing: 16) {
                    VStack(alignment: .leading, spacing: 8) {
                        Text(NSLocalizedString("scrutiny_board", comment: ""))
                            .font(.headline)
                            .fontWeight(.bold)
                            .foregroundColor(.white)
                        Text(NSLocalizedString("scrutiny_board", comment: ""))
                            .font(.subheadline)
                            .foregroundColor(.white.opacity(0.8))
                        
                        Button(action: {}) {
                            Text(NSLocalizedString("scrutiny_board", comment: ""))
                                .fontWeight(.semibold)
                                .frame(maxWidth: .infinity)
                                .padding()
                                .background(Color.orange)
                                .foregroundColor(.white)
                                .cornerRadius(10)
                        }
                    }
                    .padding()
                    .background(Color(red: 0.12, green: 0.16, blue: 0.23))
                    .cornerRadius(16)
                }
                .padding()
            }
            .navigationTitle(NSLocalizedString("scrutiny_board", comment: ""))
        }
    }
}
