import SwiftUI

public struct ParentCanteenMenuModel: Identifiable, Equatable {
    public let id: String
    public let title: String
    public let firstCourse: String
    public let secondCourse: String
    public let dessertAndDrink: String
    public let statusText: String

    public init(id: String = UUID().uuidString, title: String, firstCourse: String, secondCourse: String, dessertAndDrink: String, statusText: String = "Consumato") {
        self.id = id
        self.title = title
        self.firstCourse = firstCourse
        self.secondCourse = secondCourse
        self.dessertAndDrink = dessertAndDrink
        self.statusText = statusText
    }
}

public struct ParentCanteenView: View {
    public var menus: [ParentCanteenMenuModel]

    public init(menus: [ParentCanteenMenuModel] = []) {
        self.menus = menus
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("parent_dashboard_title", comment: ""))) {
                    if menus.isEmpty {
                        Text(NSLocalizedString("parent_dashboard_title", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(menus) { menu in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(menu.title)
                                        .fontWeight(.bold)
                                    Spacer()
                                    Text(menu.statusText)
                                        .font(.caption2)
                                        .padding(.horizontal, 6)
                                        .padding(.vertical, 2)
                                        .background(Color.green.opacity(0.2))
                                        .foregroundColor(.green)
                                        .cornerRadius(4)
                                }
                                Text(menu.firstCourse)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(menu.secondCourse)
                                    .font(.subheadline)
                                    .foregroundColor(.blue)
                                Text(menu.dessertAndDrink)
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
