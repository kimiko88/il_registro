import SwiftUI

struct StudentCreditsView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Riepilogo Crediti del Triennio (D.Lgs. 62/2017)")) {
                    HStack {
                        Text("Classe 3ª (Media 7.6)")
                        Spacer()
                        Text("11 / 12 Punti")
                            .fontWeight(.bold)
                            .foregroundColor(.blue)
                    }

                    HStack {
                        Text("Classe 4ª (Media 8.2)")
                        Spacer()
                        Text("12 / 13 Punti")
                            .fontWeight(.bold)
                            .foregroundColor(.blue)
                    }

                    HStack {
                        Text("Totale Credito Accumulato")
                            .fontWeight(.bold)
                        Spacer()
                        Text("23 Punti")
                            .font(.headline)
                            .foregroundColor(.green)
                    }
                }
            }
            .navigationTitle("Credito Scolastico")
        }
    }
}
