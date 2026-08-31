import SwiftUI

struct ParentAbsenceMonitoringView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Monte Ore e Limite Assenze (D.P.R. 122/09)")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Assenze Effettuate")
                                .fontWeight(.bold)
                            Spacer()
                            Text("54 Ore / 247 Max")
                                .font(.caption2)
                                .foregroundColor(.green)
                        }
                        ProgressView(value: 0.22)
                            .padding(.vertical, 2)
                        Text("Stato: Regolare (Anno Pienamente Valido)")
                            .font(.caption)
                            .foregroundColor(.green)
                    }
                    .padding(.vertical, 4)

                    HStack {
                        Text("Monte Ore Annuale")
                        Spacer()
                        Text("990 Ore")
                            .fontWeight(.bold)
                    }
                }
            }
            .navigationTitle("Monitoraggio Assenze")
        }
    }
}
