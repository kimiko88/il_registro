import SwiftUI

public struct GlobalSearchResultModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let categoryAndDate: String

    public init(id: String = UUID().uuidString, title: String, categoryAndDate: String) {
        self.id = id
        self.title = title
        self.categoryAndDate = categoryAndDate
    }
}

public struct GlobalSearchView: View {
    @State private var searchText = ""
    public var results: [GlobalSearchResultModel]
    public var onSearch: ((String) -> Void)?

    public init(results: [GlobalSearchResultModel] = [], onSearch: ((String) -> Void)? = nil) {
        self.results = results
        self.onSearch = onSearch
    }

    public var body: some View {
        NavigationView {
            List {
                if results.isEmpty {
                    Text(NSLocalizedString("dashboard_title", comment: ""))
                        .font(.caption)
                        .foregroundColor(.secondary)
                } else {
                    Section(header: Text(NSLocalizedString("dashboard_title", comment: ""))) {
                        ForEach(results) { item in
                            VStack(alignment: .leading, spacing: 4) {
                                Text(item.title)
                                    .fontWeight(.bold)
                                Text(item.categoryAndDate)
                                    .font(.caption)
                                    .foregroundColor(.blue)
                            }
                            .padding(.vertical, 4)
                        }
                    }
                }
            }
            .searchable(text: $searchText, prompt: NSLocalizedString("dashboard_title", comment: ""))
            .onChange(of: searchText) { newValue in
                onSearch?(newValue)
            }
            .navigationTitle(NSLocalizedString("dashboard_title", comment: ""))
        }
    }
}
