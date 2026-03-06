import type { ReactNode } from 'react';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';

const features = [
  {
    title: 'Zero Configuration',
    description: 'No YAML files. No containers. Run one binary and your entire logging stack is live in under a second.',
  },
  {
    title: 'SQL Query Interface',
    description: 'Query structured logs with standard SQL directly from the terminal. Filter, group, and aggregate.',
  },
  {
    title: 'Real-Time Streaming',
    description: 'Stream live logs with litelog tail. Filter by service or level as events arrive.',
  },
  {
    title: 'Terminal Dashboard',
    description: 'A live, full-screen terminal dashboard showing ingestion rates, error counts, and top services.',
  },
  {
    title: 'Async Ingestion Pipeline',
    description: 'HTTP handler returns immediately. Logs are batched and flushed asynchronously via background goroutines.',
  },
  {
    title: 'Log Retention',
    description: 'Automatic retention policies keep your database lean. Set it once with --retention 7d and forget it.',
  },
];

export default function Home(): ReactNode {
  const { siteConfig } = useDocusaurusContext();
  return (
    <Layout title="LiteLog — Centralized logging without the infrastructure" description="A single binary log aggregation system for developers. No containers, no config. Just fast, structured logging.">

      {/* Hero */}
      <section className="hero-section">
        <div className="container">
          <div className="hero-badge">Open Source Developer Tool</div>

          <div className="hero-logo-wrapper" style={{ justifyContent: 'center' }}>
            <img src="img/logo.png" alt="LiteLog" className="hero-logo" />
            <span className="hero-wordmark">LiteLog</span>
          </div>

          <h1 className="hero-title">
            Centralized logging<br />
            <span>without the infrastructure.</span>
          </h1>

          <p className="hero-subtitle">
            A single Go binary that ingests, stores, and queries logs at high
            throughput — without Elasticsearch, Logstash, or Kibana.
          </p>

          <div className="hero-actions">
            <Link className="btn-primary" to="/docs/quick-start">
              Get Started
            </Link>
            <Link className="btn-secondary" href="https://github.com/yashnaiduu/Litelog">
              ★ Star on GitHub
            </Link>
          </div>

          <div className="hero-code">
            <div className="hero-code-header">
              <span className="dot dot-red" />
              <span className="dot dot-yellow" />
              <span className="dot dot-green" />
            </div>
            <pre>
              <span className="code-comment"># Install</span>{'\n'}
              <span className="code-cmd">git clone</span>{' '}https://github.com/yashnaiduu/Litelog{'\n'}
              <span className="code-cmd">go build</span>{' '}-o litelog cmd/litelog/main.go{'\n'}
              {'\n'}
              <span className="code-comment"># Start the server</span>{'\n'}
              <span className="code-cmd">./litelog start</span>{' '}--retention 7d{'\n'}
              {'\n'}
              <span className="code-comment"># Query your logs</span>{'\n'}
              <span className="code-cmd">./litelog query</span>{' '}"SELECT * FROM logs WHERE level='ERROR'"
            </pre>
          </div>
        </div>
      </section>

      {/* Stats */}
      <div className="stats-bar">
        <div className="container">
          <div className="stats-grid">
            <div>
              <div className="stat-value">~40MB</div>
              <div className="stat-label">RAM Usage</div>
            </div>
            <div>
              <div className="stat-value">&lt; 1s</div>
              <div className="stat-label">Startup Time</div>
            </div>
            <div>
              <div className="stat-value">Single Binary</div>
              <div className="stat-label">No Dependencies</div>
            </div>
          </div>
        </div>
      </div>

      {/* Features */}
      <section className="features-section">
        <div className="container">
          <div className="section-label">Features</div>
          <h2 className="section-title">Built for developers who ship fast</h2>
          <p className="section-subtitle">
            Everything you need for structured logging without the overhead of enterprise observability stacks.
          </p>
          <div className="features-grid">
            {features.map((f) => (
              <div key={f.title} className="feature-card">
                <div className="feature-marker" />
                <div className="feature-title">{f.title}</div>
                <p className="feature-desc">{f.description}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Benchmarks */}
      <section className="benchmark-section">
        <div className="container">
          <div className="section-label">Performance</div>
          <h2 className="section-title">Built to be lightweight</h2>
          <p className="section-subtitle">
            LiteLog is designed for low-resource environments. Compare it to popular alternatives.
          </p>
          <div className="benchmark-grid">
            <div className="benchmark-card">
              <div className="benchmark-card-header">RAM Usage</div>
              <div className="benchmark-row"><span className="benchmark-tool">ELK Stack</span><span className="benchmark-val">2 GB+</span></div>
              <div className="benchmark-row"><span className="benchmark-tool">Prometheus + Grafana</span><span className="benchmark-val">500 MB+</span></div>
              <div className="benchmark-row highlight"><span className="benchmark-tool">LiteLog</span><span className="benchmark-val">~40 MB</span></div>
            </div>
            <div className="benchmark-card">
              <div className="benchmark-card-header">Startup Time</div>
              <div className="benchmark-row"><span className="benchmark-tool">Elasticsearch</span><span className="benchmark-val">~30s</span></div>
              <div className="benchmark-row"><span className="benchmark-tool">Logstash</span><span className="benchmark-val">~15s</span></div>
              <div className="benchmark-row highlight"><span className="benchmark-tool">LiteLog</span><span className="benchmark-val">&lt; 1s</span></div>
            </div>
          </div>
        </div>
      </section>

      {/* CTA */}
      <section className="cta-section">
        <div className="container">
          <h2 className="cta-title">Ready to replace your logging stack?</h2>
          <p className="cta-subtitle">One binary. Zero configuration. Full SQL power.</p>
          <div className="hero-actions">
            <Link className="btn-primary" to="/docs/quick-start">
              Read the Docs
            </Link>
            <Link className="btn-secondary" href="https://github.com/yashnaiduu/Litelog">
              ★ Star on GitHub
            </Link>
          </div>
        </div>
      </section>

    </Layout>
  );
}
