import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import type { ReactNode } from 'react';
import styles from './styles.module.css';

const links = [
    {
        heading: 'Product',
        items: [
            { label: 'Introduction', to: '/docs/intro' },
            { label: 'Quick Start', to: '/docs/quick-start' },
            { label: 'Commands', to: '/docs/commands' },
            { label: 'Architecture', to: '/docs/architecture' },
        ],
    },
    {
        heading: 'Resources',
        items: [
            { label: 'GitHub', href: 'https://github.com/yashnaiduu/Litelog' },
            { label: 'Releases', href: 'https://github.com/yashnaiduu/Litelog/releases' },
            { label: 'Issues', href: 'https://github.com/yashnaiduu/Litelog/issues' },
            { label: 'Pull Requests', href: 'https://github.com/yashnaiduu/Litelog/pulls' },
        ],
    },
    {
        heading: 'Legal',
        items: [
            { label: 'License (MIT)', href: 'https://github.com/yashnaiduu/Litelog/blob/main/LICENSE' },
        ],
    },
];

export default function Footer(): ReactNode {
    const { siteConfig } = useDocusaurusContext();

    return (
        <footer className={styles.footer}>
            <div className={styles.inner}>

                {/* Top row */}
                <div className={styles.top}>

                    {/* Brand */}
                    <div className={styles.brand}>
                        <div className={styles.brandLogo}>
                            <img src="/Litelog/img/logo.png" alt="LiteLog" className={styles.logo} />
                            <span className={styles.wordmark}>LiteLog</span>
                        </div>
                        <p className={styles.brandTagline}>
                            Centralized logging without the infrastructure.
                        </p>
                    </div>

                    {/* Link columns */}
                    <div className={styles.columns}>
                        {links.map((col) => (
                            <div key={col.heading} className={styles.column}>
                                <div className={styles.colHeading}>{col.heading}</div>
                                <ul className={styles.colList}>
                                    {col.items.map((item) => (
                                        <li key={item.label}>
                                            {'to' in item ? (
                                                <Link className={styles.colLink} to={item.to}>{item.label}</Link>
                                            ) : (
                                                <a className={styles.colLink} href={item.href} target="_blank" rel="noopener noreferrer">{item.label}</a>
                                            )}
                                        </li>
                                    ))}
                                </ul>
                            </div>
                        ))}
                    </div>
                </div>

                {/* Bottom bar */}
                <div className={styles.bottom}>
                    <span className={styles.copyright}>
                        Copyright {new Date().getFullYear()} LiteLog. Open source under MIT license.
                    </span>
                    <div className={styles.bottomLinks}>
                        <a className={styles.bottomLink} href="https://github.com/yashnaiduu/Litelog" target="_blank" rel="noopener noreferrer">GitHub</a>
                        <a className={styles.bottomLink} href="https://github.com/yashnaiduu/Litelog/issues" target="_blank" rel="noopener noreferrer">Issues</a>
                        <a className={styles.bottomLink} href="https://github.com/yashnaiduu/Litelog/blob/main/LICENSE" target="_blank" rel="noopener noreferrer">MIT License</a>
                    </div>
                </div>
            </div>
        </footer>
    );
}
