import SwiftUI

public struct SecretarySidiFlowModel: Identifiable, Equatable {
    public let id: String
    public let flowName: String
    public let description: String

    public init(id: String = UUID().uuidString, flowName: String, description: String = "") {
        self.id = id
        self.flowName = flowName
        self.description = description
    }
}

public struct SecretarySidiSyncView: View {
    public var flows: [SecretarySidiFlowModel]
    public var onExportFlow: ((String) -> Void)?

    public init(flows: [SecretarySidiFlowModel] = [], onExportFlow: ((String) -> Void)? = nil) {
        self.flows = flows
        self.onExportFlow = onExportFlow
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("secretary_dashboard_title", comment: ""))) {
                    if flows.isEmpty {
                        Text(NSLocalizedString("secretary_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(flows) { item in
                            HStack {
                                VStack(alignment: .leading, spacing: 2) {
                                    Text(item.flowName)
                                        .fontWeight(.bold)
                                    if !item.description.isEmpty {
                                        Text(item.description)
                                            .font(.caption)
                                            .foregroundColor(.secondary)
                                    }
                                }
                                Spacer()
                                Button("XML") {
                                    onExportFlow?(item.id)
                                }
                                .buttonStyle(.bordered)
                            }
                        }
                    }
                }
            }
            .navigationTitle(NSLocalizedString("secretary_dashboard_title", comment: ""))
        }
    }
}
