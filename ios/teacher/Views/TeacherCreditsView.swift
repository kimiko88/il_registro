import SwiftUI

struct TeacherCreditsView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Proposte Crediti - Classe 3ª A")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Mario Rossi (Media 7.60)")
                                .fontWeight(.bold)
                            Spacer()
                            Text("11 Punti")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("Banda 10-11 • Fascia max per PCTO e crediti")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Crediti Scolastici")
        }
    }
}
