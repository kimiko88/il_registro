import SwiftUI

public struct ParentQueueModel: Identifiable, Equatable {
    public let id: String
    public let teacherAndSubject: String
    public let queuePosition: String
    public let estimatedTimeAndLocation: String
    public let waitTimeDetails: String

    public init(id: String = UUID().uuidString, teacherAndSubject: String, queuePosition: String, estimatedTimeAndLocation: String, waitTimeDetails: String) {
        self.id = id
        self.teacherAndSubject = teacherAndSubject
        self.queuePosition = queuePosition
        self.estimatedTimeAndLocation = estimatedTimeAndLocation
        self.waitTimeDetails = waitTimeDetails
    }
}

public struct ParentGeneralMeetingsView: View {
    public var queues: [ParentQueueModel]

    public init(queues: [ParentQueueModel] = []) {
        self.queues = queues
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("parent_dashboard_title", comment: ""))) {
                    if queues.isEmpty {
                        Text(NSLocalizedString("parent_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(queues) { q in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(q.teacherAndSubject)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(q.queuePosition)
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(Color.green.opacity(0.2))
                                        .foregroundColor(.green)
                                        .cornerRadius(4)
                                }
                                Text(q.estimatedTimeAndLocation)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(q.waitTimeDetails)
                                    .font(.caption)
                                    .foregroundColor(.secondary)
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
