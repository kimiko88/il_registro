import SwiftUI

public struct StudentCreditModel: Identifiable, Equatable {
    public let id: String
    public let yearDescription: String
    public let pointsText: String
    public let isTotal: Bool

    public init(id: String = UUID().uuidString, yearDescription: String, pointsText: String, isTotal: Bool = false) {
        self.id = id
        self.yearDescription = yearDescription
        self.pointsText = pointsText
        self.isTotal = isTotal
    }
}

public struct StudentCreditsView: View {
    public var credits: [StudentCreditModel]

    public init(credits: [StudentCreditModel] = []) {
        self.credits = credits
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("dashboard_title", comment: ""))) {
                    if credits.isEmpty {
                        Text(NSLocalizedString("dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(credits) { credit in
                            HStack {
                                Text(credit.yearDescription)
                                    .fontWeight(credit.isTotal ? .bold : .regular)
                                Spacer()
                                Text(credit.pointsText)
                                    .fontWeight(.bold)
                                    .foregroundColor(credit.isTotal ? .green : .blue)
                            }
                        }
                    }
                }
            }
            .navigationTitle(NSLocalizedString("dashboard_title", comment: ""))
        }
    }
}
