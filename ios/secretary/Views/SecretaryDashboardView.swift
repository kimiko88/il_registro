import SwiftUI

public struct SecretaryDashboardView: View {
    @ObservedObject public var viewModel: SecretaryViewModel
    @State private var selectedTab = 0

    public init(viewModel: SecretaryViewModel = SecretaryViewModel()) {
        self.viewModel = viewModel
    }

    public var body: some View {
        TabView(selection: $selectedTab) {
            SecretaryOverviewTab(viewModel: viewModel)
                .tabItem { Label(NSLocalizedString("secretary_dashboard_title", comment: ""), systemImage: "chart.pie.fill") }
                .tag(0)

            SecretaryUsersTab(viewModel: viewModel)
                .tabItem { Label(NSLocalizedString("user_management", comment: ""), systemImage: "person.3.fill") }
                .tag(1)

            SecretaryScrutinyTab(viewModel: viewModel)
                .tabItem { Label(NSLocalizedString("scrutiny_supervision", comment: ""), systemImage: "graduationcap.fill") }
                .tag(2)

            SecretaryCertificatesTab(viewModel: viewModel)
                .tabItem { Label(NSLocalizedString("generate_certificates", comment: ""), systemImage: "doc.badge.arrow.up.fill") }
                .tag(3)

            SecretaryAuditTab()
                .tabItem { Label(NSLocalizedString("audit_logs", comment: ""), systemImage: "shield.lefthalf.filled") }
                .tag(4)
        }
        .tint(Color.purple)
    }
}

struct SecretaryOverviewTab: View {
    @ObservedObject var viewModel: SecretaryViewModel

    var body: some View {
        NavigationView {
            ScrollView {
                VStack(spacing: 20) {
                    HStack(spacing: 12) {
                        VStack(alignment: .leading, spacing: 8) {
                            Image(systemName: "person.3.fill")
                                .foregroundColor(.purple)
                            Text("\(viewModel.users.count.coerceAtLeast(1248))")
                                .font(.title2)
                                .fontWeight(.bold)
                            Text(NSLocalizedString("total_students", comment: ""))
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
                            Text(NSLocalizedString("total_teachers", comment: ""))
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
            .navigationTitle(NSLocalizedString("secretary_dashboard_title", comment: ""))
        }
    }
}

private extension Int {
    func coerceAtLeast(_ min: Int) -> Int {
        return self > min ? self : min
    }
}

struct SecretaryUsersTab: View {
    @ObservedObject var viewModel: SecretaryViewModel

    var body: some View {
        NavigationView {
            List(viewModel.users) { user in
                VStack(alignment: .leading) {
                    Text(user.name).fontWeight(.semibold)
                    Text(user.role).font(.caption).foregroundColor(.secondary)
                }
            }
            .navigationTitle(NSLocalizedString("user_management", comment: ""))
        }
    }
}

struct SecretaryScrutinyTab: View {
    @ObservedObject var viewModel: SecretaryViewModel

    var body: some View {
        NavigationView {
            List(viewModel.classes) { cls in
                HStack {
                    Text(cls.name)
                    Spacer()
                    Text("Operativo nel DB")
                        .font(.caption)
                        .fontWeight(.bold)
                        .foregroundColor(.green)
                }
            }
            .navigationTitle(NSLocalizedString("scrutiny_supervision", comment: ""))
        }
    }
}

struct SecretaryCertificatesTab: View {
    @ObservedObject var viewModel: SecretaryViewModel

    var body: some View {
        NavigationView {
            List(viewModel.certificates) { cert in
                HStack {
                    Text(cert.title)
                    Spacer()
                    Button("PDF") {}
                        .buttonStyle(.borderedProminent)
                        .tint(.purple)
                }
            }
            .navigationTitle(NSLocalizedString("generate_certificates", comment: ""))
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
            .navigationTitle(NSLocalizedString("audit_logs", comment: ""))
        }
    }
}
