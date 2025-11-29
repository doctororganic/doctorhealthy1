import './globals.css'
import { Navigation } from '@/components/layout/Navigation'

export const metadata = {
  title: 'Nutrition Platform',
  description: 'Comprehensive nutrition and fitness tracking platform',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body className="bg-gray-50">
        <Navigation />
        <main>{children}</main>
      </body>
    </html>
  )
}
