import SwiftUI

struct StudentGradeSimulatorView: View {
    @State private var simulatedGrade = "8.0"

    var body: some View {
        NavigationView {
            Form {
                Section(header: Text("Simulazione Matematica (Attuale: 7.2)")) {
                    TextField("Voto ipotetico", text: $simulatedGrade)

                    HStack {
                        Text("Nuova Media Stimata")
                            .fontWeight(.bold)
                        Spacer()
                        Text("7.45")
                            .font(.headline)
                            .foregroundColor(.blue)
                    }
                }
            }
            .navigationTitle("Simulatore Voti")
        }
    }
}
