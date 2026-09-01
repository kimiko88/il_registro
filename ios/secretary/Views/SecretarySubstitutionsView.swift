import SwiftUI

public struct SecretarySubstitutionModel: Identifiable, Equatable {
    public let id: String
    public let classAndHour: String
    public let proposedTeacher: String
    public let statusText: String

    public init(id: String = UUID().uuidString, classAndHour: String, proposedTeacher: String, statusText: String = "Docente Assente") {
        self.id = id
        self.classAndHour = classAndHour
        self.proposedTeacher = proposedTeacher
        self.statusText = statusText
    }
}

public struct SecretarySubstitutionsView: View {
    public var substitutions: [SecretarySubstitutionModel]
    public var onAssign: ((String) -> Void)?

    public init(substitutions: [SecretarySubstitutionModel] = [], onAssign: ((String) -> Void)? = nil) {
        self.substitutions = substitutions
        self.onAssign = onAssign
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("secretary_dashboard_title", comment: ""))) {
                    if substitutions.isEmpty {
                        Text(NSLocalizedString("secretary_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(substitutions) { item in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(item.classAndHour)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(item.statusText)
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(Color.red.opacity(0.2))
                                        .foregroundColor(.red)
                                        .cornerRadius(4)
                                }
                                Text(item.proposedTeacher)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)

                                Button(action: { onAssign?(item.id) }) {
                                    Text(NSLocalizedString("secretary_dashboard_title", comment: ""))
                                }
                                .buttonStyle(.borderedProminent)
                                .padding(.top, 2)
                            }
                            .padding(.vertical, 4)
                        }
                    }
                }
            }
            .navigationTitle(NSLocalizedString("secretary_dashboard_title", comment: ""))
        }
    }
}
