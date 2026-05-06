export default function TermsPage() {
  return (
    <div style={{ maxWidth: '800px', margin: '0 auto', padding: '2rem' }}>
      <h1>Terms of Service</h1>
      <p>Last updated: {new Date().toLocaleDateString()}</p>
      
      <section>
        <h2>1. Acceptance of Terms</h2>
        <p>By accessing and using this service, you accept and agree to be bound by the terms and provision of this agreement.</p>
      </section>
      
      <section>
        <h2>2. Use License</h2>
        <p>Permission is granted to temporarily download one copy of the materials for personal, non-commercial transitory viewing only.</p>
      </section>
      
      <section>
        <h2>3. Disclaimer</h2>
        <p>The materials on this website are provided on an 'as is' basis. We make no warranties, expressed or implied, and hereby disclaim and negate all other warranties including, without limitation, implied warranties or conditions of merchantability, fitness for a particular purpose, or non-infringement of intellectual property or other violation of rights.</p>
      </section>
      
      <section>
        <h2>4. Limitations</h2>
        <p>In no event shall our company or its suppliers be liable for any damages (including, without limitation, damages for loss of data or profit, or due to business interruption) arising out of the use or inability to use the materials on our website.</p>
      </section>
      
      <section>
        <h2>5. Contact Information</h2>
        <p>If you have any questions about these Terms, please contact us.</p>
      </section>
    </div>
  );
}