import SwiftUI

public struct StudentCapolavoroModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let schoolYear: String
    public let description: String
    public let selfEvaluation: String

    public init(id: String = UUID().uuidString, title: String, schoolYear: String, description: String, selfEvaluation: String) {
        self.id = id
        self.title = title
        self.schoolYear = schoolYear
        self.description = description
        self.selfEvaluation = selfEvaluation
    }
}

public struct StudentEPortfolioView: View {
    public var capolavori: [StudentCapolavoroModel]
    public var onAddCapolavoro: ((String, String, String) -> Void)?

    @State private var showingAddSheet = false
    @State private var title = ""
    @State private var description = ""
    @State private var reflection = ""

    public init(capolavori: [StudentCapolavoroModel] = [], onAddCapolavoro: ((String, String, String) -> Void)? = nil) {
        self.capolavori = capolavori
        self.onAddCapolavoro = onAddCapolavoro
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("dashboard_title", comment: ""))) {
                    if capolavori.isEmpty {
                        Text(NSLocalizedString("dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(capolavori) { item in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(item.title)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(item.schoolYear)
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(Color.purple.opacity(0.15))
                                        .foregroundColor(.purple)
                                        .cornerRadius(4)
                                }
                                Text(item.description)
                                    .font(.subheadline)
                                    .foregroundColor(.secondary)
                                Text(item.selfEvaluation)
                                    .font(.caption)
                                    .foregroundColor(.primary)
                                    .padding(6)
                                    .background(Color.gray.opacity(0.1))
                                    .cornerRadius(6)
                            }
                            .padding(.vertical, 4)
                        }
                    }
                }
            }
            .navigationTitle(NSLocalizedString("dashboard_title", comment: ""))
            .toolbar {
                Button(action: { showingAddSheet = true }) {
                    Image(systemName: "plus")
                }
            }
            .sheet(isPresented: $showingAddSheet) {
                NavigationView {
                    Form {
                        Section(header: Text(NSLocalizedString("details", comment: ""))) {
                            TextField(NSLocalizedString("title", comment: ""), text: $title)
                            TextField(NSLocalizedString("description", comment: ""), text: $description)
                            TextField(NSLocalizedString("reflection", comment: ""), text: $reflection)
                        }
                    }
                    .navigationTitle(NSLocalizedString("add", comment: ""))
                    .toolbar {
                        ToolbarItem(placement: .cancellationAction) {
                            Button(NSLocalizedString("cancel", comment: "")) { showingAddSheet = false }
                        }
                        ToolbarItem(placement: .confirmationAction) {
                            Button(NSLocalizedString("save", comment: "")) {
                                onAddCapolavoro?(title, description, reflection)
                                title = ""
                                description = ""
                                reflection = ""
                                showingAddSheet = false
                            }
                        }
                    }
                }
            }
        }
    }
}
