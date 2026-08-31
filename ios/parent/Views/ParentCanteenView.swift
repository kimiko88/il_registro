import SwiftUI

struct ParentCanteenView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Menu Mensa del Giorno")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Menu Primaverile")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Consumato")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("Primo: Risotto con verdure BIO")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Secondo: Filetto di platessa al forno")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Frutta: Mela biologica italiana")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Mensa Scolastica")
        }
    }
}
