import SwiftUI

struct StudentOrientationHoursView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Moduli Orientamento 30h (D.M. 328/22)")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Monte Ore 30/30h")
                                .fontWeight(.bold)
                            Spacer()
                            Text("100% Completato")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        ProgressView(value: 1.0)
                            .padding(.vertical, 2)
                        Text("Modulo 1: Soft Skills (10h) • Convalidato")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Modulo 2: Open Day & Lavoro (10h) • Convalidato")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Modulo 3: Capolavoro E-Portfolio (10h) • Convalidato")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Orientamento 30h")
        }
    }
}
