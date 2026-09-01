import SwiftUI

public struct SecretaryCircularModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let targetAndDate: String
    public let signatureDetails: String
    public let readRateText: String

    public init(id: String = UUID().uuidString, title: String, targetAndDate: String, signatureDetails: String, readRateText: String = "89% Letti") {
        self.id = id
        self.title = title
        self.targetAndDate = targetAndDate
        self.signatureDetails = signatureDetails
        self.readRateText = readRateText
    }
}

public struct SecretaryCommunicationsView: View {
    public var circulars: [SecretaryCircularModel]

    public init(circulars: [SecretaryCircularModel] = []) {
        self.circulars = circulars
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("secretary_dashboard_title", comment: ""))) {
                    if circulars.isEmpty {
                        Text(NSLocalizedString("secretary_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(circulars) { item in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(item.title)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(item.readRateText)
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(Color.green.opacity(0.2))
                                        .foregroundColor(.green)
                                        .cornerRadius(4)
                                }
                                Text(item.targetAndDate)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(item.signatureDetails)
                                    .font(.caption)
                                    .foregroundColor(.secondary)
                            }
                            .padding(.vertical, 4)
                        }
                    }
                }
            }
            .navigationTitle(NSLocalizedString("secretary_dashboard_title", comment: ""))
        }
    }
}
