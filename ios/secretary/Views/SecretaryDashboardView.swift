import SwiftUI

public struct SecretaryDashboardView: View {
    @ObservedObject public var viewModel: SecretaryViewModel
    public var token: String
    @State private var selectedTab = 0
    public var auditLogs: [String] = []

    public init(token: String = "", viewModel: SecretaryViewModel = SecretaryViewModel(), auditLogs: [String] = []) {
        self.token = token
        self.viewModel = viewModel
        self.auditLogs = auditLogs
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

            SecretaryAuditTab(auditLogs: auditLogs)
                .tabItem { Label(NSLocalizedString("audit_logs", comment: ""), systemImage: "shield.lefthalf.filled") }
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
                            Text("\(viewModel.users.count)")
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
                            Text("\(viewModel.classes.count)")
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

struct SecretaryUsersTab: View {
    @ObservedObject var viewModel: SecretaryViewModel

    var body: some View {
        NavigationView {
            List {
                if viewModel.users.isEmpty {
                    Text(NSLocalizedString("user_management", comment: ""))
                        .font(.caption)
                        .foregroundColor(.secondary)
                } else {
                    ForEach(viewModel.users) { user in
                        VStack(alignment: .leading) {
                            Text(user.name).fontWeight(.semibold)
                            Text(user.role).font(.caption).foregroundColor(.secondary)
                        }
                    }
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
            List {
                if viewModel.classes.isEmpty {
                    Text(NSLocalizedString("scrutiny_supervision", comment: ""))
                        .font(.caption)
                        .foregroundColor(.secondary)
                } else {
                    ForEach(viewModel.classes) { cls in
                        HStack {
                            Text(cls.name)
                            Spacer()
                            Text("Operativo")
                                .font(.caption)
                                .fontWeight(.bold)
                                .foregroundColor(.green)
                        }
                    }
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
            List {
                if viewModel.certificates.isEmpty {
                    Text(NSLocalizedString("generate_certificates", comment: ""))
                        .font(.caption)
                        .foregroundColor(.secondary)
                } else {
                    ForEach(viewModel.certificates) { cert in
                        HStack {
                            Text(cert.title)
                            Spacer()
                            Button("PDF") {}
                                .buttonStyle(.borderedProminent)
                                .tint(.purple)
                        }
                    }
                }
            }
            .navigationTitle(NSLocalizedString("generate_certificates", comment: ""))
        }
    }
}

struct SecretaryAuditTab: View {
    var auditLogs: [String] = []

    var body: some View {
        NavigationView {
            List {
                if auditLogs.isEmpty {
                    Text(NSLocalizedString("audit_logs", comment: ""))
                        .font(.caption)
                        .foregroundColor(.secondary)
                } else {
                    ForEach(auditLogs, id: \.self) { log in
                        Label(log, systemImage: "shield.fill")
                            .font(.caption)
                    }
                }
            }
            .navigationTitle(NSLocalizedString("audit_logs", comment: ""))
        }
    }
}
