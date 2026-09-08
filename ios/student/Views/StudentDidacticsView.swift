import SwiftUI

public struct DidacticMaterialModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let subjectAndTeacher: String
    public let fileDetails: String

    public init(id: String = UUID().uuidString, title: String, subjectAndTeacher: String, fileDetails: String) {
        self.id = id
        self.title = title
        self.subjectAndTeacher = subjectAndTeacher
        self.fileDetails = fileDetails
    }
}

public struct StudentDidacticsView: View {
    public var materials: [DidacticMaterialModel]
    public var onDownload: ((String) -> Void)?

    public init(materials: [DidacticMaterialModel] = [], onDownload: ((String) -> Void)? = nil) {
        self.materials = materials
        self.onDownload = onDownload
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("dashboard_title", comment: ""))) {
                    if materials.isEmpty {
                        Text(NSLocalizedString("dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(materials) { mat in
                            HStack {
                                VStack(alignment: .leading, spacing: 4) {
                                    Text(mat.title)
                                        .fontWeight(.bold)
                                    Text("\(mat.subjectAndTeacher) • \(mat.fileDetails)")
                                        .font(.caption)
                                        .foregroundColor(.secondary)
                                }
                                Spacer()
                                Button(action: { onDownload?(mat.id) }) {
                                    Image(systemName: "arrow.down.circle.fill")
                                        .font(.title2)
                                        .foregroundColor(.blue)
                                }
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
