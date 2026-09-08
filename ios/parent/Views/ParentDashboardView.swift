import SwiftUI

public struct ParentGradeDisplayModel: Identifiable, Equatable {
    public let id: String
    public let subject: String
    public let gradeWithDetails: String

    public init(id: String = UUID().uuidString, subject: String, gradeWithDetails: String) {
        self.id = id
        self.subject = subject
        self.gradeWithDetails = gradeWithDetails
    }
}

public struct ParentDashboardView: View {
    @ObservedObject public var viewModel: ParentViewModel
    public var token: String
    @State private var selectedTab = 0
    public var grades: [ParentGradeDisplayModel] = []
    public var colloquioSlots: [ParentColloquioSlotModel] = []
    public var circulars: [String] = []

    public init(
        token: String = "",
        viewModel: ParentViewModel = ParentViewModel(),
        grades: [ParentGradeDisplayModel] = [],
        colloquioSlots: [ParentColloquioSlotModel] = [],
        circulars: [String] = []
    ) {
        self.token = token
        self.viewModel = viewModel
        self.grades = grades
        self.colloquioSlots = colloquioSlots
        self.circulars = circulars
    }

    public var body: some View {
        TabView(selection: $selectedTab) {
            ParentHomeTab(viewModel: viewModel)
                .tabItem { Label(NSLocalizedString("parent_dashboard_title", comment: ""), systemImage: "house.fill") }
                .tag(0)

            ParentGradesTab(grades: grades)
                .tabItem { Label(NSLocalizedString("child_grades", comment: ""), systemImage: "chart.bar.doc.horizontal.fill") }
                .tag(1)

            ParentAttendanceTab(viewModel: viewModel)
                .tabItem { Label(NSLocalizedString("pending_justifications", comment: ""), systemImage: "checkmark.circle.fill") }
                .tag(2)

            ParentColloquiTab(colloquioSlots: colloquioSlots)
                .tabItem { Label(NSLocalizedString("upcoming_colloqui", comment: ""), systemImage: "calendar.badge.clock") }
                .tag(3)

            ParentCircularsTab(circulars: circulars)
                .tabItem { Label(NSLocalizedString("communications", comment: ""), systemImage: "megaphone.fill") }
                .tag(4)
        }
        .tint(Color(red: 0.05, green: 0.58, blue: 0.53))
        .task {
            if !token.isEmpty {
                await viewModel.loadFromDatabase(token: token)
            }
        }
    }
}

struct ParentHomeTab: View {
    @ObservedObject var viewModel: ParentViewModel

    var body: some View {
        NavigationView {
            ScrollView {
                VStack(spacing: 20) {
                    if viewModel.children.isEmpty {
                        Text(NSLocalizedString("no_children_found", comment: ""))
                            .font(.subheadline)
                            .foregroundColor(.secondary)
                            .padding(.top, 40)
                    } else {
                        ForEach(viewModel.children) { child in
                            VStack(alignment: .leading, spacing: 12) {
                                Text("\(child.firstName) \(child.lastName)")
                                    .font(.title2)
                                    .fontWeight(.bold)
                                    .foregroundColor(.white)
                                Text(child.className)
                                    .font(.subheadline)
                                    .foregroundColor(.white.opacity(0.8))
                                HStack(spacing: 16) {
                                    Text(NSLocalizedString("parent_dashboard_title", comment: ""))
                                        .font(.caption)
                                        .fontWeight(.semibold)
                                        .padding(8)
                                        .background(Color.white.opacity(0.2))
                                        .cornerRadius(8)
                                        .foregroundColor(.white)
                                }
                            }
                            .frame(maxWidth: .infinity, alignment: .leading)
                            .padding()
                            .background(Color(red: 0.06, green: 0.15, blue: 0.28))
                            .cornerRadius(16)
                        }
                    }
                }
                .padding()
            }
            .navigationTitle(NSLocalizedString("parent_dashboard_title", comment: ""))
        }
    }
}

struct ParentGradesTab: View {
    var grades: [ParentGradeDisplayModel]

    var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("child_grades", comment: ""))) {
                    if grades.isEmpty {
                        Text(NSLocalizedString("no_grades", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(grades) { item in
                            HStack {
                                Text(item.subject)
                                Spacer()
                                Text(item.gradeWithDetails).foregroundColor(.green).fontWeight(.bold)
                            }
                        }
                    }
                }
            }
            .navigationTitle(NSLocalizedString("child_grades", comment: ""))
        }
    }
}

struct ParentAttendanceTab: View {
    @ObservedObject var viewModel: ParentViewModel

    var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("pending_justifications", comment: ""))) {
                    if viewModel.absences.isEmpty {
                        Text(NSLocalizedString("no_pending_justifications", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(viewModel.absences) { item in
                            HStack {
                                VStack(alignment: .leading) {
                                    Text(item.date)
                                        .fontWeight(.bold)
                                    Text(item.type)
                                        .font(.caption)
                                        .foregroundColor(.secondary)
                                }
                                Spacer()
                                if item.isJustified {
                                    Text(NSLocalizedString("status_justified", comment: ""))
                                        .font(.caption)
                                        .foregroundColor(.green)
                                } else {
                                    Button(NSLocalizedString("btn_justify", comment: "")) {
                                        _ = viewModel.justifyAbsence(id: item.id, reason: "Motivata")
                                    }
                                    .buttonStyle(.borderedProminent)
                                    .tint(Color(red: 0.05, green: 0.58, blue: 0.53))
                                }
                            }
                        }
                    }
                }
            }
            .navigationTitle(NSLocalizedString("pending_justifications", comment: ""))
        }
    }
}

struct ParentColloquiTab: View {
    var colloquioSlots: [ParentColloquioSlotModel]

    var body: some View {
        NavigationView {
            List {
                if colloquioSlots.isEmpty {
                    Text(NSLocalizedString("no_data", comment: ""))
                        .font(.caption)
                        .foregroundColor(.secondary)
                } else {
                    ForEach(colloquioSlots) { slot in
                        HStack {
                            VStack(alignment: .leading) {
                                Text(slot.teacherName)
                                    .fontWeight(.semibold)
                                Text(slot.timeSlot)
                                    .font(.caption)
                                    .foregroundColor(.secondary)
                            }
                            Spacer()
                            Button(NSLocalizedString("book_colloquio", comment: "")) {}
                                .buttonStyle(.borderedProminent)
                                .tint(Color(red: 0.05, green: 0.58, blue: 0.53))
                        }
                    }
                }
            }
            .navigationTitle(NSLocalizedString("upcoming_colloqui", comment: ""))
        }
    }
}

struct ParentCircularsTab: View {
    var circulars: [String]

    var body: some View {
        NavigationView {
            List {
                if circulars.isEmpty {
                    Text(NSLocalizedString("no_circulars", comment: ""))
                        .font(.caption)
                        .foregroundColor(.secondary)
                } else {
                    ForEach(circulars, id: \.self) { c in
                        Text(c)
                    }
                }
            }
            .navigationTitle(NSLocalizedString("communications", comment: ""))
        }
    }
}
