import SwiftUI

public struct TeacherGradeWeightModel: Identifiable, Equatable {
    public let id: String
    public let typeName: String
    public let weightPercent: Int

    public init(id: String = UUID().uuidString, typeName: String, weightPercent: Int) {
        self.id = id
        self.typeName = typeName
        self.weightPercent = weightPercent
    }
}

public struct TeacherGradeWeightsView: View {
    public var weights: [TeacherGradeWeightModel]

    public init(weights: [TeacherGradeWeightModel] = []) {
        self.weights = weights
    }

    public var body: some View {
        NavigationView {
            List {
                Section(header: Text(NSLocalizedString("teacher_dashboard_title", comment: ""))) {
                    if weights.isEmpty {
                        Text(NSLocalizedString("add_grade", comment: ""))
                            .font(.caption)
                            .foregroundColor(.secondary)
                    } else {
                        ForEach(weights) { item in
                            HStack {
                                Text(item.typeName)
                                Spacer()
                                Text("\(item.weightPercent)%")
                                    .fontWeight(.bold)
                                    .foregroundColor(.blue)
                            }
                        }
                    }
                }
            }
            .navigationTitle(NSLocalizedString("teacher_dashboard_title", comment: ""))
        }
    }
}
