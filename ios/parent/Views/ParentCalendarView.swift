import SwiftUI

struct ParentCalendarView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Festività e Chiusure")) {
                    VStack(alignment: .leading, spacing: 4) {
                        HStack {
                            Text("Ponte Festa della Repubblica")
                                .fontWeight(.bold)
                            Spacer()
                            Text("1-2 Giugno")
                                .font(.caption2)
                                .foregroundColor(.blue)
                        }
                        Text("Sospensione delle lezioni deliberate dal Consiglio")
                            .font(.subheadline)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)

                    VStack(alignment: .leading, spacing: 4) {
                        HStack {
                            Text("Termine Lezioni A.S. 2025/2026")
                                .fontWeight(.bold)
                            Spacer()
                            Text("6 Giugno")
                                .font(.caption2)
                                .foregroundColor(.green)
                        }
                        Text("Ultimo giorno di scuola • Uscita ore 12:00")
                            .font(.subheadline)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Calendario Scolastico")
        }
    }
}
