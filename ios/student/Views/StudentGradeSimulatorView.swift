import SwiftUI

public struct StudentGradeSimulatorView: View {
    public var subjectName: String
    public var currentAverage: Double
    @State private var simulatedGrade = "8.0"

    public init(subjectName: String = "Matematica", currentAverage: Double = 7.2) {
        self.subjectName = subjectName
        self.currentAverage = currentAverage
    }

    public var calculatedAverage: Double {
        let gradeVal = Double(simulatedGrade) ?? currentAverage
        return (currentAverage * 3.0 + gradeVal) / 4.0
    }

    public var body: some View {
        NavigationView {
            Form {
                Section(header: Text("\(subjectName) (Attuale: \(String(format: "%.2f", currentAverage)))")) {
                    TextField("Voto ipotetico", text: $simulatedGrade)

                    HStack {
                        Text(NSLocalizedString("grades_title", comment: ""))
                            .fontWeight(.bold)
                        Spacer()
                        Text(String(format: "%.2f", calculatedAverage))
                            .font(.headline)
                            .foregroundColor(.blue)
                    }
                }
            }
            .navigationTitle(NSLocalizedString("grades_title", comment: ""))
        }
    }
}
