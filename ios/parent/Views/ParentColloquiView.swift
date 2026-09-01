import SwiftUI

public struct ParentColloquioSlotModel: Identifiable, Equatable {
    public let id: String
    public let teacherAndSubject: String
    public let timeAndLocation: String
    public let statusText: String

    public init(id: String = UUID().uuidString, teacherAndSubject: String, timeAndLocation: String, statusText: String = "Disponibile") {
        self.id = id
        self.teacherAndSubject = teacherAndSubject
        self.timeAndLocation = timeAndLocation
        self.statusText = statusText
    }
}

public struct ParentColloquiView: View {
    public var slots: [ParentColloquioSlotModel]
    public var onBookSlot: ((String) -> Void)?

    public init(slots: [ParentColloquioSlotModel] = [], onBookSlot: ((String) -> Void)? = nil) {
        self.slots = slots
        self.onBookSlot = onBookSlot
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("parent_dashboard_title", comment: ""))) {
                    if slots.isEmpty {
                        Text(NSLocalizedString("parent_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(slots) { slot in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(slot.teacherAndSubject)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(slot.statusText)
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(Color.green.opacity(0.2))
                                        .foregroundColor(.green)
                                        .cornerRadius(4)
                                }
                                Text(slot.timeAndLocation)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)

                                Button(NSLocalizedString("sign_hour", comment: "")) {
                                    onBookSlot?(slot.id)
                                }
                                .buttonStyle(.borderedProminent)
                                .padding(.top, 2)
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
