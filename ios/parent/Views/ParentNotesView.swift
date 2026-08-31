import SwiftUI

public struct ParentDisciplinaryNoteModel: Identifiable, Equatable {
    public let id: String
    public let teacherAndDate: String
    public let description: String
    public var isAcknowledged: Bool

    public init(id: String = UUID().uuidString, teacherAndDate: String, description: String, isAcknowledged: Bool = false) {
        self.id = id
        self.teacherAndDate = teacherAndDate
        self.description = description
        self.isAcknowledged = isAcknowledged
    }
}

public struct ParentNotesView: View {
    public var notes: [ParentDisciplinaryNoteModel]
    public var onAcknowledge: ((String) -> Void)?
    @State private var acknowledgedIds = Set<String>()

    public init(notes: [ParentDisciplinaryNoteModel] = [], onAcknowledge: ((String) -> Void)? = nil) {
        self.notes = notes
        self.onAcknowledge = onAcknowledge
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("parent_dashboard_title", comment: ""))) {
                    if notes.isEmpty {
                        Text(NSLocalizedString("parent_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(notes) { note in
                            let isAck = note.isAcknowledged || acknowledgedIds.contains(note.id)
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(NSLocalizedString("pending_justifications", comment: ""))
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(isAck ? "Firmata" : "In Attesa")
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(isAck ? Color.green.opacity(0.2) : Color.red.opacity(0.2))
                                        .foregroundColor(isAck ? .green : .red)
                                        .cornerRadius(4)
                                }
                                Text(note.teacherAndDate)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(note.description)
                                    .font(.caption)
                                    .foregroundColor(.secondary)

                                if !isAck {
                                    Button(NSLocalizedString("sign_hour", comment: "")) {
                                        acknowledgedIds.insert(note.id)
                                        onAcknowledge?(note.id)
                                    }
                                    .buttonStyle(.borderedProminent)
                                    .padding(.top, 2)
                                }
                            }
                            .padding(.vertical, 4)
                        }
                    }
                }
            }
            .navigationTitle(NSLocalizedString("parent_dashboard_title", comment: ""))
        }
    }
}
