import SwiftUI

struct TeacherGradeWeightsView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Pesi Percentuali per Tipologia di Prova")) {
                    HStack {
                        Text("Prova Scritta / Compito in Classe")
                        Spacer()
                        Text("100%")
                            .fontWeight(.bold)
                            .foregroundColor(.blue)
                    }

                    HStack {
                        Text("Interrogazione Orale")
                        Spacer()
                        Text("100%")
                            .fontWeight(.bold)
                            .foregroundColor(.blue)
                    }

                    HStack {
                        Text("Test Breve / Verifica Formativa")
                        Spacer()
                        Text("50%")
                            .fontWeight(.bold)
                            .foregroundColor(.orange)
                    }
                }
            }
            .navigationTitle("Configurazione Pesi Voti")
        }
    }
}
