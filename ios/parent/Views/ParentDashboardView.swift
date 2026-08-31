import SwiftUI

struct ChildModel: Identifiable, Hashable {
    let id = UUID()
    let name: String
    let className: String
}

struct ParentDashboardView: View {
    let children = [
        ChildModel(name: "Marco Rossi", className: "Classe 2A"),
        ChildModel(name: "Giulia Rossi", className: "Classe 4B")
    ]
    
    @State private var selectedChild: ChildModel
    @State private var selectedTab = 0
    
    init() {
        _selectedChild = State(initialValue: children.first!)
    }

    var body: some View {
        TabView(selection: $selectedTab) {
            ParentHomeTab(selectedChild: selectedChild, children: children)
                .tabItem { Label("Figli", systemImage: "figure.2.and.child.holdinghands") }
                .tag(0)

            ParentGradesTab(selectedChild: selectedChild)
                .tabItem { Label("Voti", systemImage: "chart.bar.doc.horizontal.fill") }
                .tag(1)

            ParentAttendanceTab(selectedChild: selectedChild)
                .tabItem { Label("Giustifiche", systemImage: "checkmark.circle.fill") }
                .tag(2)

            ParentColloquiTab()
                .tabItem { Label("Colloqui", systemImage: "calendar.badge.clock") }
                .tag(3)

            ParentCircularsTab()
                .tabItem { Label("Circolari", systemImage: "megaphone.fill") }
                .tag(4)
        }
        .tint(Color(red: 0.05, green: 0.58, blue: 0.53))
    }
}

struct ParentHomeTab: View {
    let selectedChild: ChildModel
    let children: [ChildModel]

    var body: some View {
        NavigationView {
            ScrollView {
                VStack(spacing: 20) {
                    VStack(alignment: .leading, spacing: 12) {
                        Text(selectedChild.name)
                            .font(.title2)
                            .fontWeight(.bold)
                            .foregroundColor(.white)
                        Text(selectedChild.className)
                            .font(.subheadline)
                            .foregroundColor(.white.opacity(0.8))
                        HStack(spacing: 16) {
                            Text("Media Generale: 8.1")
                                .font(.caption)
                                .fontWeight(.semibold)
                                .padding(8)
                                .background(Color.white.opacity(0.2))
                                .cornerRadius(8)
                                .foregroundColor(.white)
                            Text("Assenze Totali: 2")
                                .font(.caption)
                                .fontWeight(.semibold)
                                .padding(8)
                                .background(Color.white.opacity(0.2))
                                .cornerRadius(8)
                                .foregroundColor(.white)
                        }
                    }
                    .frame(maxWidth: .infinity, alignment: .leading)
                    .padding()
                    .background(Color(red: 0.06, green: 0.15, blue: 0.28))
                    .cornerRadius(16)
                }
                .padding()
            }
            .navigationTitle("Supervisione Figli")
        }
    }
}

struct ParentGradesTab: View {
    let selectedChild: ChildModel

    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Valutazioni 1° Quadrimestre")) {
                    HStack {
                        Text("Matematica")
                        Spacer()
                        Text("8.5 (Scritto)").foregroundColor(.green).fontWeight(.bold)
                    }
                    HStack {
                        Text("Italiano")
                        Spacer()
                        Text("8.0 (Orale)").foregroundColor(.green).fontWeight(.bold)
                    }
                    HStack {
                        Text("Fisica")
                        Spacer()
                        Text("7.5 (Scritto)").foregroundColor(.green).fontWeight(.bold)
                    }
                }
            }
            .navigationTitle("Voti Figlio")
        }
    }
}

struct ParentAttendanceTab: View {
    let selectedChild: ChildModel

    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Assenze & Ritardi da Giustificare")) {
                    HStack {
                        VStack(alignment: .leading) {
                            Text("26 Ago 2026")
                                .fontWeight(.bold)
                            Text("Assenza per motivi di salute")
                                .font(.caption)
                                .foregroundColor(.secondary)
                        }
                        Spacer()
                        Button("Giustifica") {}
                            .buttonStyle(.borderedProminent)
                            .tint(Color(red: 0.05, green: 0.58, blue: 0.53))
                    }
                }
            }
            .navigationTitle("Giustifiche")
        }
    }
}

struct ParentColloquiTab: View {
    var body: some View {
        NavigationView {
            List {
                HStack {
                    VStack(alignment: .leading) {
                        Text("Prof. Bianchi (Matematica)")
                            .fontWeight(.semibold)
                        Text("Giovedì 15:30 - 15:45")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    Spacer()
                    Button("Prenota") {}
                        .buttonStyle(.borderedProminent)
                        .tint(Color(red: 0.05, green: 0.58, blue: 0.53))
                }
            }
            .navigationTitle("Ricevimento Famiglie")
        }
    }
}

struct ParentCircularsTab: View {
    var body: some View {
        NavigationView {
            List {
                Text("Circolare n. 42: Calendario incontri scuola-famiglia")
                Text("Circolare n. 41: Protocollo viaggi di istruzione")
                Text("Circolare n. 40: Attivazione registro elettronico")
            }
            .navigationTitle("Circolari & Avvisi")
        }
    }
}
