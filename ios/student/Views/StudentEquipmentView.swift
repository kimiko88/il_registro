import SwiftUI

public struct StudentEquipmentModel: Identifiable, Equatable {
    public let id: String
    public let resourceName: String
    public let details: String
    public let isHighlighted: Bool

    public init(id: String = UUID().uuidString, resourceName: String, details: String, isHighlighted: Bool = false) {
        self.id = id
        self.resourceName = resourceName
        self.details = details
        self.isHighlighted = isHighlighted
    }
}

public struct StudentEquipmentView: View {
    public var equipment: [StudentEquipmentModel]

    public init(equipment: [StudentEquipmentModel] = []) {
        self.equipment = equipment
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("dashboard_title", comment: ""))) {
                    if equipment.isEmpty {
                        Text(NSLocalizedString("dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(equipment) { item in
                            HStack {
                                Text(item.resourceName)
                                Spacer()
                                Text(item.details)
                                    .fontWeight(.bold)
                                    .foregroundColor(item.isHighlighted ? .blue : .primary)
                            }
                        }
                    }
                }
            }
            .navigationTitle(NSLocalizedString("dashboard_title", comment: ""))
        }
    }
}
