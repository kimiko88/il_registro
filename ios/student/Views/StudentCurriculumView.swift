import SwiftUI

public struct CurriculumSectionModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let description: String

    public init(id: String = UUID().uuidString, title: String, description: String) {
        self.id = id
        self.title = title
        self.description = description
    }
}

public struct StudentCurriculumView: View {
    public var sections: [CurriculumSectionModel]
    public var onDownloadPdf: (() -> Void)?

    public init(sections: [CurriculumSectionModel] = [], onDownloadPdf: (() -> Void)? = nil) {
        self.sections = sections
        self.onDownloadPdf = onDownloadPdf
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("dashboard_title", comment: ""))) {
                    HStack {
                        VStack(alignment: .leading, spacing: 4) {
                            Text(NSLocalizedString("dashboard_title", comment: ""))
                                .fontWeight(.bold)
                            Text(NSLocalizedString("grades_title", comment: ""))
                                .font(.caption)
                                .foregroundColor(.secondary)
                        }
                        Spacer()
                        Button(action: { onDownloadPdf?() }) {
                            Label("PDF", systemImage: "arrow.down.doc.fill")
                        }
                        .buttonStyle(.borderedProminent)
                    }
                }

                if !sections.isEmpty {
                    Section(header: Text("Sezioni")) {
                        ForEach(sections) { sec in
                            VStack(alignment: .leading, spacing: 4) {
                                Text(sec.title)
                                    .font(.headline)
                                Text(sec.description)
                                    .font(.subheadline)
                                    .foregroundColor(.secondary)
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
