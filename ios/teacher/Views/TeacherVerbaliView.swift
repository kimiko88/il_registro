import SwiftUI

public struct TeacherVerbaleModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let sessionDetails: String
    public let statusText: String

    public init(id: String = UUID().uuidString, title: String, sessionDetails: String, statusText: String = "Approvato") {
        self.id = id
        self.title = title
        self.sessionDetails = sessionDetails
        self.statusText = statusText
    }
}

public struct TeacherVerbaliView: View {
    public var verbali: [TeacherVerbaleModel]
    public var onViewPdf: ((String) -> Void)?

    public init(verbali: [TeacherVerbaleModel] = [], onViewPdf: ((String) -> Void)? = nil) {
        self.verbali = verbali
        self.onViewPdf = onViewPdf
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("teacher_dashboard_title", comment: ""))) {
                    if verbali.isEmpty {
                        Text(NSLocalizedString("teacher_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(verbali) { item in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(item.title)
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
                                Text(item.sessionDetails)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)

                                Button(action: { onViewPdf?(item.id) }) {
                                    Label("PDF", systemImage: "doc.text.fill")
                                }
                                .buttonStyle(.bordered)
                                .padding(.top, 2)
                            }
                            .padding(.vertical, 4)
                        }
                    }
                }
            }
            .navigationTitle(NSLocalizedString("teacher_dashboard_title", comment: ""))
        }
    }
}
