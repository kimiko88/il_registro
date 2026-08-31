import SwiftUI

struct StudentEquipmentView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Dotazioni e Risorse Assegnate")) {
                    HStack {
                        Text("Armadietto Personale")
                        Spacer()
                        Text("ARM-104 (Piano 1)")
                            .fontWeight(.bold)
                            .foregroundColor(.blue)
                    }

                    HStack {
                        Text("Badge Accesso NFC")
                        Spacer()
                        Text("NFC-99120-STUD")
                            .fontWeight(.bold)
                    }

                    HStack {
                        Text("Chromebook Comodato")
                        Spacer()
                        Text("ASUS 14\" (CB-8821)")
                            .fontWeight(.bold)
                            .foregroundColor(.green)
                    }
                }
            }
            .navigationTitle("Dotazioni & Badge")
        }
    }
}
