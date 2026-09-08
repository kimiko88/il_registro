import SwiftUI

public struct SecretaryClassFormationModel: Identifiable, Equatable {
    public let id: String
    public let classTitle: String
    public let studentComposition: String
    public let distributionDetails: String
    public let scoreBadgeText: String

    public init(id: String = UUID().uuidString, classTitle: String, studentComposition: String, distributionDetails: String, scoreBadgeText: String = "Ottimale 98%") {
        self.id = id
        self.classTitle = classTitle
        self.studentComposition = studentComposition
        self.distributionDetails = distributionDetails
        self.scoreBadgeText = scoreBadgeText
    }
}

public struct SecretaryClassFormationView: View {
    public var classes: [SecretaryClassFormationModel]

    public init(classes: [SecretaryClassFormationModel] = []) {
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
                                    Text(item.classTitle)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(item.scoreBadgeText)
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(Color.green.opacity(0.2))
                                        .foregroundColor(.green)
                                        .cornerRadius(4)
                                }
                                Text(item.studentComposition)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(item.distributionDetails)
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
