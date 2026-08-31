import SwiftUI

struct SecretaryFascicoloView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Fascicolo Studente: Mario Rossi (3ª A)")) {
                    HStack {
                        Text("Codice Fiscale")
                        Spacer()
                        Text("RSSMRA08A01H501Z")
                            .fontWeight(.bold)
                    }

                    HStack {
                        Text("Codice SIDI")
                        Spacer()
                        Text("SIDI-1049281")
                            .fontWeight(.bold)
                            .foregroundColor(.blue)
                    }

                    HStack {
                        Text("Stato Vaccinale")
                        Spacer()
                        Text("Regolare")
                            .fontWeight(.bold)
                            .foregroundColor(.green)
                    }
                }
            }
            .navigationTitle("Fascicolo Studente")
        }
    }
}
