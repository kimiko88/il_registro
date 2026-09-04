#if canImport(UIKit)
@_exported import UIKit
#elseif canImport(AppKit)
import AppKit
import SwiftUI

public struct UIColor {
    public static let systemBackground = NSColor.windowBackgroundColor
    public static let secondarySystemGroupedBackground = NSColor.controlBackgroundColor
    public static let systemGroupedBackground = NSColor.windowBackgroundColor
}

public extension Color {
    init(_ nsColor: NSColor) {
        self.init(nsColor: nsColor)
    }
}

public enum UITextAutocapitalizationType {
    case none
    case words
    case sentences
    case allCharacters
}

public enum UIKeyboardType {
    case `default`
    case asciiCapable
    case numbersAndPunctuation
    case URL
    case numberPad
    case phonePad
    case namePhonePad
    case emailAddress
    case decimalPad
    case twitter
    case webSearch
    case asciiCapableNumberPad
}

public extension View {
    func autocapitalization(_ type: UITextAutocapitalizationType) -> some View {
        self
    }

    func keyboardType(_ type: UIKeyboardType) -> some View {
        self
    }
}
#endif
