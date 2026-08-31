import SwiftUI

public struct TeacherUdaModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let subjects: String
    public let targetCompetencies: String
    public let statusText: String

    public init(id: String = UUID().uuidString, title: String, subjects: String, targetCompetencies: String, statusText: String = "In Corso") {
        self.id = id
        self.title = title
        self.subjects = subjects
        self.targetCompetencies = targetCompetencies
        self.statusText = statusText
    }
}

public struct TeacherUdaView: View {
    public var udas: [TeacherUdaModel]

    public init(udas: [TeacherUdaModel] = []) {
        self.udas = udas
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("teacher_dashboard_title", comment: ""))) {
                    if udas.isEmpty {
                        Text(NSLocalizedString("teacher_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(udas) { item in
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
                                Text(item.subjects)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(item.targetCompetencies)
                                    .font(.caption)
                                    .foregroundColor(.secondary)
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
