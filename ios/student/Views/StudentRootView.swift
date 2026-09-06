import SwiftUI

public struct StudentRootView: View {
    @State private var isLoggedIn: Bool
    @State private var token: String
    @State private var studentName: String

    public init(startLoggedIn: Bool = false, initialToken: String = "", initialStudentName: String = "Studente") {
        self._isLoggedIn = State(initialValue: startLoggedIn)
        self._token = State(initialValue: initialToken)
        self._studentName = State(initialValue: initialStudentName)
    }

    public var body: some View {
        if isLoggedIn {
            StudentDashboardView(token: token, studentName: studentName) {
                withAnimation {
                    self.isLoggedIn = false
                    self.token = ""
                    self.studentName = "Studente"
                }
            }
        } else {
            LoginView(isLoggedIn: $isLoggedIn, token: $token, studentName: $studentName)
        }
    }
}
