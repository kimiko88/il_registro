import SwiftUI

struct StudentSchoolCalendarView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Festività e Sospensioni")) {
                    VStack(alignment: .leading, spacing: 4) {
                        HStack {
                            Text("Festa della Repubblica & Ponte")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Ponte Deliberato")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.blue.opacity(0.2))
                                .foregroundColor(.blue)
                                .cornerRadius(4)
                        }
                        Text("1 - 2 Giugno 2026")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Delibera N. 42 del Consiglio d'Istituto")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)

                    VStack(alignment: .leading, spacing: 4) {
                        HStack {
                            Text("Termine delle Lezioni")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Ultimo Giorno")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("6 Giugno 2026")
                            .font(.subheadline)
                            .foregroundColor(.green)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Calendario Scolastico")
        }
    }
}
