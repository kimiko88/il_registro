import SwiftUI

struct SecretaryDashboardView: View {
    @State private var selectedTab = 0

    var body: some View {
        TabView(selection: $selectedTab) {
            SecretaryOverviewTab()
                .tabItem { Label("Pannello", systemImage: "chart.pie.fill") }
                .tag(0)

            SecretaryUsersTab()
                .tabItem { Label("Utenti", systemImage: "person.3.fill") }
                .tag(1)

            SecretaryScrutinyTab()
                .tabItem { Label("Scrutini", systemImage: "graduationcap.fill") }
                .tag(2)

            SecretaryCertificatesTab()
                .tabItem { Label("Certificati", systemImage: "doc.badge.arrow.up.fill") }
                .tag(3)

            SecretaryAuditTab()
                .tabItem { Label("Audit", systemImage: "shield.lefthalf.filled") }
                .tag(4)
        }
        .tint(Color.purple)
    }
}

struct SecretaryOverviewTab: View {
    var body: some View {
        NavigationView {
            ScrollView {
                VStack(spacing: 20) {
                    HStack(spacing: 12) {
                        VStack(alignment: .leading, spacing: 8) {
                            Image(systemName: "person.3.fill")
                                .foregroundColor(.purple)
                            Text("1,248")
                                .font(.title2)
                                .fontWeight(.bold)
                            Text("Studenti Iscritti")
                                .font(.caption)
                                .foregroundColor(.secondary)
                        }
                        .padding()
                        .frame(maxWidth: .infinity, alignment: .leading)
                        .background(Color(UIColor.secondarySystemGroupedBackground))
                        .cornerRadius(14)

                        VStack(alignment: .leading, spacing: 8) {
                            Image(systemName: "graduationcap.fill")
                                .foregroundColor(.purple)
                            Text("94")
                                .font(.title2)
                                .fontWeight(.bold)
                            Text("Docenti Attivi")
                                .font(.caption)
                                .foregroundColor(.secondary)
                        }
                        .padding()
                        .frame(maxWidth: .infinity, alignment: .leading)
                        .background(Color(UIColor.secondarySystemGroupedBackground))
                        .cornerRadius(14)
                    }
                }
                .padding()
            }
            .navigationTitle("Pannello Segreteria")
        }
    }
}

struct SecretaryUsersTab: View {
    let users = [
        ("Prof.ssa Maria Rossi", "Docente • Lettere"),
        ("Prof. Marco Bianchi", "Docente • Matematica"),
        ("Mario Rossi (2B)", "Studente"),
        ("Giuseppe Rossi", "Genitore")
    ]

    var body: some View {
        NavigationView {
            List(users, id: \.0) { user, role in
                VStack(alignment: .leading) {
                    Text(user).fontWeight(.semibold)
                    Text(role).font(.caption).foregroundColor(.secondary)
                }
            }
            .navigationTitle("Anagrafica Utenti")
        }
    }
}

struct SecretaryScrutinyTab: View {
    let classes = [("1A", "Completato"), ("2A", "In corso"), ("3A", "Completato"), ("4B", "Differito")]

    var body: some View {
        NavigationView {
            List(classes, id: \.0) { cls, status in
                HStack {
                    Text("Classe \(cls)")
                    Spacer()
                    Text(status)
                        .font(.caption)
                        .fontWeight(.bold)
                        .foregroundColor(status == "Completato" ? .green : .orange)
                }
            }
            .navigationTitle("Supervisione Scrutini")
        }
    }
}

struct SecretaryCertificatesTab: View {
    let certs = [
        "Certificato di Iscrizione e Frequenza",
        "Certificato con Valutazioni e Voti",
        "Certificato di Diploma"
    ]

    var body: some View {
        NavigationView {
            List(certs, id: \.self) { cert in
                HStack {
                    Text(cert)
                    Spacer()
                    Button("PDF") {}
                        .buttonStyle(.borderedProminent)
                        .tint(.purple)
                }
            }
            .navigationTitle("Emissione Certificati")
        }
    }
}

struct SecretaryAuditTab: View {
    let logs = [
        "31 Ago 11:20 • Accesso SuperAdmin",
        "31 Ago 10:45 • Generazione Certificato Matricola 412",
        "31 Ago 09:30 • Chiusura Scrutinio Classe 3A"
    ]

    var body: some View {
        NavigationView {
            List(logs, id: \.self) { log in
                Label(log, systemImage: "shield.fill")
                    .font(.caption)
            }
            .navigationTitle("Registro Audit")
        }
    }
}
