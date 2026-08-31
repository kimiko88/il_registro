import SwiftUI

struct TeacherRegisterView: View {
    @State private var selectedTab = 0
    @State private var lessonTopic = "Equazioni e disequazioni esponenziali"
    @State private var isHourSigned = true

    var body: some View {
        TabView(selection: $selectedTab) {
            TeacherFirmaTab(topic: $lessonTopic, isSigned: $isHourSigned)
                .tabItem { Label("Registro", systemImage: "pencil.and.list.clipboard") }
                .tag(0)

            TeacherAppelloTab()
                .tabItem { Label("Appello", systemImage: "checkmark.rectangle.stack.fill") }
                .tag(1)

            TeacherGradesTab()
                .tabItem { Label("Voti", systemImage: "chart.bar.fill") }
                .tag(2)

            TeacherAgendaTab()
                .tabItem { Label("Agenda", systemImage: "calendar.badge.plus") }
                .tag(3)

            TeacherScrutinyTab()
                .tabItem { Label("Scrutini", systemImage: "graduationcap.fill") }
                .tag(4)
        }
        .tint(Color.blue)
    }
}

struct TeacherFirmaTab: View {
    @Binding var topic: String
    @Binding var isSigned: Bool

    var body: some View {
        NavigationView {
            ScrollView {
                VStack(spacing: 20) {
                    VStack(alignment: .leading, spacing: 12) {
                        Text("Firma Registro di Classe")
                            .font(.headline)
                            .fontWeight(.bold)
                        Text("Classe 3A • Matematica (1a e 2a ora)")
                            .font(.subheadline)
                            .foregroundColor(.secondary)

                        TextField("Argomento della Lezione", text: $topic)
                            .textFieldStyle(.roundedBorder)

                        Button(action: { isSigned = true }) {
                            HStack {
                                Image(systemName: isSigned ? "checkmark.circle.fill" : "pencil.line")
                                Text(isSigned ? "Ora Lezione Firmata" : "Firma Ora Lezione")
                                    .fontWeight(.semibold)
                            }
                            .frame(maxWidth: .infinity)
                            .padding()
                            .background(isSigned ? Color.green : Color.blue)
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
            .navigationTitle("Registro di Classe")
        }
    }
}

struct TeacherAppelloTab: View {
    let students = ["Banchi Andrea", "Bianchi Elena", "Ferrari Matteo", "Rossi Sofia"]

    var body: some View {
        NavigationView {
            List(students, id: \.self) { student in
                HStack {
                    Text(student).fontWeight(.medium)
                    Spacer()
                    Button("Presente") {}
                        .buttonStyle(.borderedProminent)
                        .tint(.green)
                }
            }
            .navigationTitle("Appello & Presenze")
        }
    }
}

struct TeacherGradesTab: View {
    var body: some View {
        NavigationView {
            List {
                HStack {
                    Text("Banchi Andrea")
                    Spacer()
                    Text("8 (Scritto)").foregroundColor(.blue).fontWeight(.bold)
                }
                HStack {
                    Text("Bianchi Elena")
                    Spacer()
                    Text("7.5 (Orale)").foregroundColor(.blue).fontWeight(.bold)
                }
            }
            .navigationTitle("Inserimento Voti")
        }
    }
}

struct TeacherAgendaTab: View {
    var body: some View {
        NavigationView {
            List {
                Button(action: {}) {
                    Label("Assegna Nuovo Compito in Agenda", systemImage: "plus.circle.fill")
                }
            }
            .navigationTitle("Agenda di Classe")
        }
    }
}

struct TeacherScrutinyTab: View {
    var body: some View {
        NavigationView {
            ScrollView {
                VStack(spacing: 16) {
                    VStack(alignment: .leading, spacing: 8) {
                        Text("Tabellone Scrutini & Scrutini Differiti")
                            .font(.headline)
                            .fontWeight(.bold)
                            .foregroundColor(.white)
                        Text("Delibera voti finali, condotta e verifiche debiti")
                            .font(.subheadline)
                            .foregroundColor(.white.opacity(0.8))
                        
                        Button(action: {}) {
                            Text("Sessione Scrutinio Differito (Debiti)")
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
            .navigationTitle("Gestione Scrutini")
        }
    }
}
