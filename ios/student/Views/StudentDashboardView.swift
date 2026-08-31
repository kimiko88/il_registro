import SwiftUI

struct GradeModel: Identifiable {
    let id = UUID()
    let subject: String
    let value: Double
    let date: String
    let type: String
}

struct HomeworkModel: Identifiable {
    let id = UUID()
    let subject: String
    let description: String
    let dueDate: String
    var isDone: Bool
}

struct StudentDashboardView: View {
    @State private var studentName: String = "Mario Rossi"
    @State private var gpaAverage: Double = 7.8
    @State private var selectedTab = 0

    var body: some View {
        TabView(selection: $selectedTab) {
            StudentHomeView(studentName: studentName, gpaAverage: gpaAverage)
                .tabItem {
                    Label("Home", systemImage: "house.fill")
                }
                .tag(0)

            StudentGradesView()
                .tabItem {
                    Label("Voti", systemImage: "chart.bar.doc.horizontal.fill")
                }
                .tag(1)

            StudentAgendaView()
                .tabItem {
                    Label("Agenda", systemImage: "calendar")
                }
                .tag(2)

            StudentAttendanceView()
                .tabItem {
                    Label("Presenze", systemImage: "checkmark.circle.fill")
                }
                .tag(3)

            StudentReportCardView()
                .tabItem {
                    Label("Pagella", systemImage: "doc.text.fill")
                }
                .tag(4)
        }
        .tint(Color.purple)
    }
}

// 1. Home View
struct StudentHomeView: View {
    let studentName: String
    let gpaAverage: Double

    let recentGrades = [
        GradeModel(subject: "Matematica", value: 8.5, date: "30 Ago", type: "Scritto"),
        GradeModel(subject: "Italiano", value: 7.5, date: "28 Ago", type: "Orale"),
        GradeModel(subject: "Inglese", value: 8.0, date: "25 Ago", type: "Pratico"),
        GradeModel(subject: "Fisica", value: 7.0, date: "22 Ago", type: "Scritto")
    ]

    var body: some View {
        NavigationView {
            ScrollView {
                VStack(spacing: 20) {
                    // Youthful Hero Header Card
                    VStack(alignment: .leading, spacing: 12) {
                        HStack {
                            VStack(alignment: .leading, spacing: 4) {
                                Text("Ciao, \(studentName)!")
                                    .font(.title)
                                    .fontWeight(.bold)
                                    .foregroundColor(.white)
                                Text("Media Generale: \(String(format: "%.1f", gpaAverage))")
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

                    // Recent Grades Section
                    VStack(alignment: .leading, spacing: 12) {
                        Text("Ultimi Voti Inseriti")
                            .font(.headline)
                            .fontWeight(.bold)
                            .accessibilityAddTraits(.isHeader)

                        ForEach(recentGrades) { grade in
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
                                Text(String(format: "%.1f", grade.value))
                                    .font(.headline)
                                    .fontWeight(.bold)
                                    .foregroundColor(.white)
                                    .padding(.horizontal, 12)
                                    .padding(.vertical, 6)
                                    .background(grade.value >= 6.0 ? Color.green : Color.red)
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
            .navigationTitle("La Mia Dashboard")
            .background(Color(UIColor.systemGroupedBackground).ignoresSafeArea())
        }
    }
}

// 2. Grades View
struct StudentGradesView: View {
    let allGrades = [
        GradeModel(subject: "Matematica", value: 8.5, date: "30 Ago", type: "Scritto"),
        GradeModel(subject: "Italiano", value: 8.0, date: "28 Ago", type: "Tema"),
        GradeModel(subject: "Inglese", value: 9.0, date: "25 Ago", type: "Pratico"),
        GradeModel(subject: "Fisica", value: 7.0, date: "22 Ago", type: "Scritto"),
        GradeModel(subject: "Storia", value: 8.0, date: "18 Ago", type: "Orale"),
        GradeModel(subject: "Filosofia", value: 8.5, date: "12 Ago", type: "Orale")
    ]

    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Riepilogo Quadrimestre")) {
                    HStack {
                        Text("Media Voti Totale")
                            .fontWeight(.semibold)
                        Spacer()
                        Text("7.9")
                            .fontWeight(.bold)
                            .foregroundColor(.purple)
                    }
                }

                Section(header: Text("Tutte le Valutazioni")) {
                    ForEach(allGrades) { grade in
                        HStack {
                            VStack(alignment: .leading) {
                                Text(grade.subject)
                                    .fontWeight(.semibold)
                                Text("\(grade.type) • \(grade.date)")
                                    .font(.caption)
                                    .foregroundColor(.secondary)
                            }
                            Spacer()
                            Text(String(format: "%.1f", grade.value))
                                .fontWeight(.bold)
                                .foregroundColor(grade.value >= 6 ? .green : .red)
                        }
                    }
                }
            }
            .navigationTitle("I Miei Voti")
        }
    }
}

// 3. Agenda View
struct StudentAgendaView: View {
    @State private var tasks = [
        HomeworkModel(subject: "Matematica", description: "Esercizi pag 142 disequazioni esponenziali", dueDate: "Domani", isDone: false),
        HomeworkModel(subject: "Fisica", description: "Relazione moto rettilineo uniforme", dueDate: "Tra 2 giorni", isDone: false),
        HomeworkModel(subject: "Italiano", description: "Capitolo 8 I Promessi Sposi", dueDate: "Tra 3 giorni", isDone: true)
    ]

    var body: some View {
        NavigationView {
            List {
                ForEach($tasks) { $task in
                    HStack {
                        Image(systemName: task.isDone ? "checkmark.circle.fill" : "circle")
                            .foregroundColor(task.isDone ? .green : .gray)
                            .onTapGesture { task.isDone.toggle() }
                        VStack(alignment: .leading) {
                            Text(task.subject)
                                .fontWeight(.semibold)
                            Text(task.description)
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
            .navigationTitle("Compiti & Agenda")
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
            .navigationTitle("Presenze & Assenze")
        }
    }
}

// 5. Report Card View
struct StudentReportCardView: View {
    let reportGrades = [
        ("Italiano", 8), ("Latino", 7), ("Inglese", 8),
        ("Storia", 8), ("Matematica", 8), ("Fisica", 7),
        ("Scienze", 8), ("Arte", 9), ("Scienze Motorie", 9)
    ]

    var body: some View {
        NavigationView {
            ScrollView {
                VStack(spacing: 16) {
                    VStack(alignment: .leading, spacing: 8) {
                        Text("Documento di Valutazione Finale")
                            .font(.headline)
                            .fontWeight(.bold)
                            .foregroundColor(.white)
                        Text("Esito Scrutinio: PROMOSSO / AMMESSO")
                            .font(.subheadline)
                            .foregroundColor(.white.opacity(0.9))
                        Text("Condotta: 9 • Credito: 8")
                            .font(.caption)
                            .foregroundColor(.white.opacity(0.8))
                    }
                    .frame(maxWidth: .infinity, alignment: .leading)
                    .padding()
                    .background(Color.purple)
                    .cornerRadius(16)

                    ForEach(reportGrades, id: \.0) { subject, grade in
                        HStack {
                            Text(subject).fontWeight(.medium)
                            Spacer()
                            Text("\(grade)")
                                .fontWeight(.bold)
                                .foregroundColor(.white)
                                .padding(.horizontal, 12)
                                .padding(.vertical, 4)
                                .background(grade >= 6 ? Color.green : Color.red)
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
            .navigationTitle("La Mia Pagella")
            .background(Color(UIColor.systemGroupedBackground).ignoresSafeArea())
        }
    }
}
