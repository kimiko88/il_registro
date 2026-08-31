import SwiftUI

public struct SecretaryTextbookClassModel: Identifiable, Equatable {
    public let id: String
    public let classNameAndTrack: String
    public let costDetails: String
    public let adoptionDetails: String
    public let statusText: String

    public init(id: String = UUID().uuidString, classNameAndTrack: String, costDetails: String, adoptionDetails: String, statusText: String = "Entro il Limite") {
        self.id = id
        self.classNameAndTrack = classNameAndTrack
        self.costDetails = costDetails
        self.adoptionDetails = adoptionDetails
        self.statusText = statusText
    }
}

public struct SecretaryTextbooksView: View {
    public var classes: [SecretaryTextbookClassModel]

    public init(classes: [SecretaryTextbookClassModel] = []) {
        self.classes = classes
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("secretary_dashboard_title", comment: ""))) {
                    if classes.isEmpty {
                        Text(NSLocalizedString("secretary_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(classes) { item in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(item.classNameAndTrack)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(item.statusText)
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(Color.green.opacity(0.2))
                                        .foregroundColor(.green)
                                        .cornerRadius(4)
                                }
                                Text(item.costDetails)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(item.adoptionDetails)
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
