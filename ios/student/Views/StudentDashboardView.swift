import SwiftUI

public struct StudentDashboardView: View {
    @ObservedObject public var viewModel: StudentViewModel
    @State private var studentName: String = "Mario Rossi"
    @State private var selectedTab = 0

    public init(viewModel: StudentViewModel = StudentViewModel()) {
        self.viewModel = viewModel
    }

    public var body: some View {
        TabView(selection: $selectedTab) {
            StudentHomeView(studentName: studentName, viewModel: viewModel)
                .tabItem {
                    Label(NSLocalizedString("dashboard_title", comment: ""), systemImage: "house.fill")
                }
                .tag(0)

            StudentGradesView(viewModel: viewModel)
                .tabItem {
                    Label(NSLocalizedString("grades_title", comment: ""), systemImage: "chart.bar.doc.horizontal.fill")
                }
                .tag(1)

            StudentAgendaView(viewModel: viewModel)
                .tabItem {
                    Label(NSLocalizedString("agenda_title", comment: ""), systemImage: "calendar")
                }
                .tag(2)

            StudentAttendanceView()
                .tabItem {
                    Label(NSLocalizedString("attendance_title", comment: ""), systemImage: "checkmark.circle.fill")
                }
                .tag(3)

            StudentReportCardView(viewModel: viewModel)
                .tabItem {
                    Label(NSLocalizedString("report_card_title", comment: ""), systemImage: "doc.text.fill")
                }
                .tag(4)
        }
        .tint(Color.purple)
    }
}

// 1. Home View
struct StudentHomeView: View {
    let studentName: String
    @ObservedObject var viewModel: StudentViewModel

    var body: some View {
        NavigationView {
            ScrollView {
                VStack(spacing: 20) {
                    VStack(alignment: .leading, spacing: 12) {
                        HStack {
                            VStack(alignment: .leading, spacing: 4) {
                                Text("Ciao, \(studentName)!")
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
                                    Text("\(grade.type) • \(grade.date)")
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
                Section(header: Text("Riepilogo Database")) {
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
                    ForEach(viewModel.grades) { grade in
                        HStack {
                            VStack(alignment: .leading) {
                                Text(grade.subject)
                                    .fontWeight(.semibold)
                                Text("\(grade.type) • \(grade.date)")
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
                            Text("Scadenza: \(task.dueDate)")
                                .font(.caption2)
                                .foregroundColor(.purple)
                        }
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle(NSLocalizedString("agenda_title", comment: ""))
        }
    }
}

// 4. Attendance View
struct StudentAttendanceView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Riepilogo Presenze Anno Scolastico")) {
                    HStack {
                        Text("Giorni di Presenza")
                        Spacer()
                        Text("184").fontWeight(.bold).foregroundColor(.green)
                    }
                    HStack {
                        Text("Assenze Totali")
                        Spacer()
                        Text("3").fontWeight(.bold).foregroundColor(.red)
                    }
                    HStack {
                        Text("Ritardi")
                        Spacer()
                        Text("1").fontWeight(.bold).foregroundColor(.orange)
                    }
                }

                Section(header: Text("Storico Giustificazioni")) {
                    HStack {
                        Text("26 Ago 2026 • Assenza")
                        Spacer()
                        Label("Giustificata", systemImage: "checkmark.seal.fill").foregroundColor(.green).font(.caption)
                    }
                    HStack {
                        Text("18 Ago 2026 • Ritardo 1a ora")
                        Spacer()
                        Label("Giustificata", systemImage: "checkmark.seal.fill").foregroundColor(.green).font(.caption)
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
                        Text("Documento di Valutazione Ufficiale")
                            .font(.headline)
                            .fontWeight(.bold)
                            .foregroundColor(.white)
                        Text("Esito Scrutinio: AMMESSO / PROMOSSO")
                            .font(.subheadline)
                            .foregroundColor(.white.opacity(0.9))
                        Text("Media Voti: \(String(format: "%.1f", viewModel.calculateGPA()))")
                            .font(.caption)
                            .foregroundColor(.white.opacity(0.8))
                    }
                    .frame(maxWidth: .infinity, alignment: .leading)
                    .padding()
                    .background(Color.purple)
                    .cornerRadius(16)

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

                    Button(action: {}) {
                        HStack {
                            Image(systemName: "arrow.down.doc.fill")
                            Text("Scarica Pagella Ufficiale (PDF)")
                                .fontWeight(.semibold)
                        }
                        .frame(maxWidth: .infinity)
                        .padding()
                        .background(Color.purple)
                        .foregroundColor(.white)
                        .cornerRadius(12)
                    }
                }
                .padding()
            }
            .navigationTitle(NSLocalizedString("report_card_title", comment: ""))
            .background(Color(UIColor.systemGroupedBackground).ignoresSafeArea())
        }
    }
}
