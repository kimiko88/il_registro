import SwiftUI

public struct AccessibilityStatementView: View {
    public var rtdEmail: String
    public var complianceStatus: String

    public init(rtdEmail: String = "rtd@scuola.edu.it", complianceStatus: String = "Totalmente Conforme") {
        self.rtdEmail = rtdEmail
        self.complianceStatus = complianceStatus
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("dashboard_title", comment: ""))) {
                    VStack(alignment: .leading, spacing: 6) {
                        Text(complianceStatus)
                            .fontWeight(.bold)
                            .foregroundColor(.green)
                        Text(NSLocalizedString("grades_title", comment: ""))
                            .font(.subheadline)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }

                Section(header: Text(NSLocalizedString("dashboard_title", comment: ""))) {
                    Text("RTD: \(rtdEmail)")
                        .font(.caption)
                        .foregroundColor(.blue)
                }
            }
            .navigationTitle(NSLocalizedString("dashboard_title", comment: ""))
        }
    }
}
