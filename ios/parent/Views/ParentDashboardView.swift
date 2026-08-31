import SwiftUI

public struct ParentDashboardView: View {
    @ObservedObject public var viewModel: ParentViewModel
    @State private var selectedTab = 0

    public init(viewModel: ParentViewModel = ParentViewModel()) {
        self.viewModel = viewModel
    }

    public var body: some View {
        TabView(selection: $selectedTab) {
            ParentHomeTab(viewModel: viewModel)
                .tabItem { Label(NSLocalizedString("my_children", comment: ""), systemImage: "figure.2.and.child.holdinghands") }
                .tag(0)

            ParentGradesTab(viewModel: viewModel)
                .tabItem { Label("Voti", systemImage: "chart.bar.doc.horizontal.fill") }
                .tag(1)

            ParentAttendanceTab(viewModel: viewModel)
                .tabItem { Label(NSLocalizedString("pending_absences", comment: ""), systemImage: "checkmark.circle.fill") }
                .tag(2)

            ParentColloquiTab()
                .tabItem { Label(NSLocalizedString("book_colloqui", comment: ""), systemImage: "calendar.badge.clock") }
                .tag(3)

            ParentCircularsTab()
                .tabItem { Label("Circolari", systemImage: "megaphone.fill") }
                .tag(4)
        }
        .tint(Color(red: 0.05, green: 0.58, blue: 0.53))
    }
}

struct ParentHomeTab: View {
    @ObservedObject var viewModel: ParentViewModel

    var body: some View {
        NavigationView {
            ScrollView {
                VStack(spacing: 20) {
                    ForEach(viewModel.children) { child in
                        VStack(alignment: .leading, spacing: 12) {
                            Text("\(child.firstName) \(child.lastName)")
                                .font(.title2)
                                .fontWeight(.bold)
                                .foregroundColor(.white)
                            Text(child.className)
                                .font(.subheadline)
                                .foregroundColor(.white.opacity(0.8))
                            HStack(spacing: 16) {
                                Text("Stato: Attivo nel database")
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
                }
                .padding()
            }
            .navigationTitle(NSLocalizedString("parent_dashboard_title", comment: ""))
        }
    }
}

struct ParentGradesTab: View {
    @ObservedObject var viewModel: ParentViewModel

    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Valutazioni Figlio")) {
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
    @ObservedObject var viewModel: ParentViewModel

    var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("pending_absences", comment: ""))) {
                    ForEach(viewModel.absences) { item in
                        HStack {
                            VStack(alignment: .leading) {
                                Text(item.date)
                                    .fontWeight(.bold)
                                Text(item.type)
                                    .font(.caption)
                                    .foregroundColor(.secondary)
                            }
                            Spacer()
                            if item.isJustified {
                                Text("Giustificata")
                                    .font(.caption)
                                    .foregroundColor(.green)
                            } else {
                                Button(NSLocalizedString("justify_action", comment: "")) {
                                    _ = viewModel.justifyAbsence(id: item.id, reason: "Motivata")
                                }
                                .buttonStyle(.borderedProminent)
                                .tint(Color(red: 0.05, green: 0.58, blue: 0.53))
                            }
                        }
                    }
                }
            }
            .navigationTitle(NSLocalizedString("pending_absences", comment: ""))
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
                    Button(NSLocalizedString("book_colloqui", comment: "")) {}
                        .buttonStyle(.borderedProminent)
                        .tint(Color(red: 0.05, green: 0.58, blue: 0.53))
                }
            }
            .navigationTitle(NSLocalizedString("book_colloqui", comment: ""))
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
