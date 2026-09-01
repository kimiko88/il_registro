import SwiftUI

public struct ExamSimulationModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let contentDescription: String
    public let durationAndGrid: String
    public let tag: String

    public init(id: String = UUID().uuidString, title: String, contentDescription: String, durationAndGrid: String, tag: String = "Traccia Ufficiale") {
        self.id = id
        self.title = title
        self.contentDescription = contentDescription
        self.durationAndGrid = durationAndGrid
        self.tag = tag
    }
}

public struct StudentExamSimulationsView: View {
    public var simulations: [ExamSimulationModel]
    public var onDownloadPdf: ((String) -> Void)?

    public init(simulations: [ExamSimulationModel] = [], onDownloadPdf: ((String) -> Void)? = nil) {
        self.simulations = simulations
        self.onDownloadPdf = onDownloadPdf
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("dashboard_title", comment: ""))) {
                    if simulations.isEmpty {
                        Text(NSLocalizedString("dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(simulations) { sim in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(sim.title)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(sim.tag)
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(Color.green.opacity(0.2))
                                        .foregroundColor(.green)
                                        .cornerRadius(4)
                                }
                                Text(sim.contentDescription)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(sim.durationAndGrid)
                                    .font(.caption)
                                    .foregroundColor(.secondary)

                                Button(action: { onDownloadPdf?(sim.id) }) {
                                    Label("PDF", systemImage: "arrow.down.doc.fill")
                                }
                                .buttonStyle(.bordered)
                                .padding(.top, 2)
                            }
                            .padding(.vertical, 4)
                        }
                    }
                }
            }
            .navigationTitle(NSLocalizedString("dashboard_title", comment: ""))
        }
    }
}
