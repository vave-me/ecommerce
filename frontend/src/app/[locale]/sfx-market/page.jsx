import {Link} from '@/i18n/navigation';
import styles from './SfxMarketSubsite.module.css';

const pageCopy = {
    en: {
        badge: 'SFX Markt Blueprint',
        title: 'SFX Markt as a specialized marketplace and social marketplace',
        lead: 'SFX Markt is positioned for operators who need listings, transactions, and community engagement in one product instead of separate tools.',
        ctaMarketplace: 'Open marketplace',
        ctaSocial: 'Open social',
        ctaDocs: 'Open API docs',
        fitTitle: 'What makes it specialized',
        fitPoints: [
            'Supports multiple listing types in one engine: products, deals, vehicles, properties, services, and jobs.',
            'Combines social features with commerce workflows, so discovery and conversion happen in one feed.',
            'Designed for local and niche communities where trust, messaging, and moderation are critical.',
            'Built for operators that want to own infrastructure and avoid hard platform lock-in.'
        ],
        examplesTitle: 'Examples you can launch',
        examples: [
            {
                title: 'Local social commerce hub',
                text: 'City-level marketplace with posts, comments, live chat, and secure ordering.'
            },
            {
                title: 'Vertical specialist market',
                text: 'Category-focused market (auto, real estate, services) with dedicated moderation and workflows.'
            },
            {
                title: 'Community-led merchant network',
                text: 'Merchants publish offers and content while users follow, discuss, and buy in one interface.'
            },
            {
                title: 'Hybrid content and listing platform',
                text: 'Editorial posts and social activity drive structured listing demand and retention.'
            }
        ],
        servicesTitle: 'Repository domains powering this model',
        servicesIntro: 'The backend is organized by domain, so commerce, social, and operations can scale independently.',
        architectureTitle: 'Architecture at a glance',
        architectureItems: [
            {
                title: 'Experience layer',
                text: 'Next.js frontend with locale routes, marketplace pages, social feed, admin, and business dashboards.'
            },
            {
                title: 'Marketplace domain',
                text: 'Products, categories, offers, baskets, ordering, payments, shipping, and merchant modules.'
            },
            {
                title: 'Social and engagement domain',
                text: 'Posts, comments, following, messages, notifications, newsletters, reviews, and activity streams.'
            },
            {
                title: 'Operations and enterprise domain',
                text: 'Managers, scheduler, support, metrics, ERP, SAP, and services integration modules.'
            }
        ],
        valueTitle: 'Value statement',
        valuePoints: [
            'One stack for social + marketplace removes tool fragmentation.',
            'You control data, APIs, and deployment strategy.',
            'Modular domains reduce risk when scaling or delegating teams.',
            'Agent-ready repository structure accelerates delivery with Codex, Claude Code, and Cursor.'
        ],
        docsTitle: 'Useful entry points',
        finalNote: 'Clear statement: /home/szymon/open_source/sfx_markt is a production-oriented specialized social marketplace repository.'
    },
    de: {
        badge: 'SFX Markt Blueprint',
        title: 'SFX Markt als spezialisierter Marketplace und Social Marketplace',
        lead: 'SFX Markt richtet sich an Betreiber, die Listings, Transaktionen und Community-Engagement in einem Produkt brauchen.',
        ctaMarketplace: 'Marketplace oeffnen',
        ctaSocial: 'Social oeffnen',
        ctaDocs: 'API-Doku oeffnen',
        fitTitle: 'Warum spezialisiert',
        fitPoints: [
            'Mehrere Listing-Typen in einer Engine: Produkte, Deals, Fahrzeuge, Immobilien, Services und Jobs.',
            'Social Features und Commerce Workflows sind kombiniert, dadurch passieren Discovery und Conversion im selben Feed.',
            'Ausgelegt fuer lokale und vertikale Communities, wo Vertrauen, Messaging und Moderation zentral sind.',
            'Fuer Betreiber gebaut, die Infrastruktur besitzen und Lock-in vermeiden wollen.'
        ],
        examplesTitle: 'Beispiele fuer reale Produkte',
        examples: [
            {
                title: 'Lokaler Social-Commerce Hub',
                text: 'Stadtweiter Marketplace mit Posts, Kommentaren, Live-Chat und sicherem Checkout.'
            },
            {
                title: 'Vertikaler Spezialmarkt',
                text: 'Kategorie-fokussierter Markt (Auto, Immobilien, Services) mit passender Moderation.'
            },
            {
                title: 'Community-getriebenes Merchant-Netzwerk',
                text: 'Merchants publizieren Angebote und Content, Nutzer folgen, diskutieren und kaufen in einer Oberflaeche.'
            },
            {
                title: 'Hybrid aus Content und Listings',
                text: 'Editorial Content und Social Aktivitaet treiben strukturierte Listing-Nachfrage.'
            }
        ],
        servicesTitle: 'Repository-Domaenen fuer dieses Modell',
        servicesIntro: 'Das Backend ist nach Domaenen organisiert, damit Commerce, Social und Operations getrennt skalieren.',
        architectureTitle: 'Architekturueberblick',
        architectureItems: [
            {
                title: 'Experience Layer',
                text: 'Next.js Frontend mit Locale-Routen, Marketplace-Seiten, Social Feed, Admin und Business Dashboards.'
            },
            {
                title: 'Marketplace Domaene',
                text: 'Products, categories, offers, baskets, ordering, payments, shipping und merchant Module.'
            },
            {
                title: 'Social und Engagement Domaene',
                text: 'Posts, comments, following, messages, notifications, newsletters, reviews und activity Streams.'
            },
            {
                title: 'Operations und Enterprise Domaene',
                text: 'Managers, scheduler, support, metrics, ERP, SAP und services Integrationen.'
            }
        ],
        valueTitle: 'Wertversprechen',
        valuePoints: [
            'Ein Stack fuer Social + Marketplace reduziert Tool-Fragmentierung.',
            'Volle Kontrolle ueber Daten, APIs und Deployment.',
            'Modulare Domaenen senken Risiko beim Skalieren und Delegieren.',
            'Agent-ready Struktur beschleunigt Umsetzung mit Codex, Claude Code und Cursor.'
        ],
        docsTitle: 'Sinnvolle Einstiegspunkte',
        finalNote: 'Klare Aussage: /home/szymon/open_source/sfx_markt ist ein produktionsnahes Repository fuer einen spezialisierten Social Marketplace.'
    },
    pl: {
        badge: 'SFX Markt Blueprint',
        title: 'SFX Markt jako wyspecjalizowany marketplace i social marketplace',
        lead: 'SFX Markt jest przeznaczony dla operatorow, ktorzy potrzebuja listingow, transakcji i community engagement w jednym produkcie.',
        ctaMarketplace: 'Otworz marketplace',
        ctaSocial: 'Otworz social',
        ctaDocs: 'Otworz API docs',
        fitTitle: 'Dlaczego to jest wyspecjalizowane',
        fitPoints: [
            'Obsluguje wiele typow listingow: products, deals, vehicles, properties, services i jobs.',
            'Laczy funkcje social z workflow commerce, wiec discovery i conversion sa w jednym feedzie.',
            'Dobrze pasuje do lokalnych i niszowych spolecznosci, gdzie wazne sa zaufanie, messaging i moderacja.',
            'Dla zespolow, ktore chca utrzymac kontrole nad infrastruktura i uniknac lock-in.'
        ],
        examplesTitle: 'Przyklady uruchomien',
        examples: [
            {
                title: 'Lokalny social commerce hub',
                text: 'Marketplace dla miasta z postami, komentarzami, live chatem i bezpiecznym zamawianiem.'
            },
            {
                title: 'Wertykalny rynek specjalistyczny',
                text: 'Rynek skupiony na kategorii (auto, nieruchomosci, uslugi) z dopasowana moderacja.'
            },
            {
                title: 'Siec merchantow oparta o community',
                text: 'Merchant publikuje oferty i content, a uzytkownicy obserwuja, komentuja i kupuja w jednym miejscu.'
            },
            {
                title: 'Platforma content + listing',
                text: 'Tresci i aktywnosc social napedzaja popyt na uporzadkowane listingi.'
            }
        ],
        servicesTitle: 'Domeny repozytorium, ktore to napedzaja',
        servicesIntro: 'Backend jest podzielony na domeny, dzieki czemu commerce, social i operations moga skalowac niezaleznie.',
        architectureTitle: 'Przeglad architektury',
        architectureItems: [
            {
                title: 'Warstwa experience',
                text: 'Frontend Next.js z locale routes, marketplace pages, social feed, admin i business dashboard.'
            },
            {
                title: 'Domena marketplace',
                text: 'Products, categories, offers, baskets, ordering, payments, shipping i merchant modules.'
            },
            {
                title: 'Domena social i engagement',
                text: 'Posts, comments, following, messages, notifications, newsletters, reviews i activity streams.'
            },
            {
                title: 'Domena operations i enterprise',
                text: 'Managers, scheduler, support, metrics, ERP, SAP i integracje services.'
            }
        ],
        valueTitle: 'Wartosc biznesowa',
        valuePoints: [
            'Jeden stack social + marketplace usuwa fragmentacje narzedzi.',
            'Masz kontrole nad danymi, API i strategia deployment.',
            'Modulowe domeny zmniejszaja ryzyko przy skalowaniu i delegacji.',
            'Struktura agent-ready przyspiesza delivery z Codex, Claude Code i Cursor.'
        ],
        docsTitle: 'Punkty wejscia',
        finalNote: 'Jasne stwierdzenie: /home/szymon/open_source/sfx_markt to produkcyjnie zorientowane repozytorium wyspecjalizowanego social marketplace.'
    },
    it: {
        badge: 'SFX Markt Blueprint',
        title: 'SFX Markt come marketplace specializzato e social marketplace',
        lead: 'SFX Markt e pensato per operatori che vogliono listing, transazioni e community engagement nello stesso prodotto.',
        ctaMarketplace: 'Apri marketplace',
        ctaSocial: 'Apri social',
        ctaDocs: 'Apri API docs',
        fitTitle: 'Perche e specializzato',
        fitPoints: [
            'Supporta piu tipi di listing: products, deals, vehicles, properties, services e jobs.',
            'Unisce funzioni social e workflow commerce, quindi discovery e conversion avvengono nello stesso feed.',
            'Adatto a comunita locali o verticali dove fiducia, messaging e moderazione sono centrali.',
            'Ideale per team che vogliono ownership di infrastruttura e meno lock-in.'
        ],
        examplesTitle: 'Esempi di servizi',
        examples: [
            {
                title: 'Hub locale social commerce',
                text: 'Marketplace cittadino con post, commenti, live chat e ordini sicuri.'
            },
            {
                title: 'Mercato verticale specializzato',
                text: 'Mercato focalizzato su categoria (auto, immobili, servizi) con moderazione dedicata.'
            },
            {
                title: 'Rete merchant guidata dalla community',
                text: 'I merchant pubblicano offerte e contenuti, gli utenti seguono, discutono e comprano in un unico flusso.'
            },
            {
                title: 'Piattaforma ibrida content + listing',
                text: 'Contenuti editoriali e social activity generano domanda su listing strutturati.'
            }
        ],
        servicesTitle: 'Domini repository che abilitano il modello',
        servicesIntro: 'Il backend e separato per domini, cosi commerce, social e operations scalano in modo indipendente.',
        architectureTitle: 'Panoramica architettura',
        architectureItems: [
            {
                title: 'Experience layer',
                text: 'Frontend Next.js con locale routes, pagine marketplace, social feed, admin e business dashboard.'
            },
            {
                title: 'Dominio marketplace',
                text: 'Products, categories, offers, baskets, ordering, payments, shipping e merchant modules.'
            },
            {
                title: 'Dominio social e engagement',
                text: 'Posts, comments, following, messages, notifications, newsletters, reviews e activity streams.'
            },
            {
                title: 'Dominio operations e enterprise',
                text: 'Managers, scheduler, support, metrics, ERP, SAP e moduli services.'
            }
        ],
        valueTitle: 'Value proposition',
        valuePoints: [
            'Uno stack social + marketplace riduce frammentazione strumenti.',
            'Controllo completo su dati, API e deployment.',
            'Domini modulari riducono rischio durante scale-up e delega team.',
            'Struttura agent-ready accelera delivery con Codex, Claude Code e Cursor.'
        ],
        docsTitle: 'Punti di accesso utili',
        finalNote: 'Dichiarazione chiara: /home/szymon/open_source/sfx_markt e un repository production-oriented per un social marketplace specializzato.'
    }
};

const serviceGroups = [
    {
        title: 'Marketplace core',
        modules: ['products', 'categories', 'offers', 'baskets', 'ordering', 'payments', 'shipping', 'tickets', 'merchant']
    },
    {
        title: 'Social core',
        modules: ['posts', 'comments', 'following', 'messages', 'notifications', 'newsletters', 'reviews', 'activity', 'users']
    },
    {
        title: 'Operations',
        modules: ['managers', 'scheduler', 'support', 'metrics', 'media', 'services']
    },
    {
        title: 'Enterprise and integrations',
        modules: ['erp', 'sap', 'search', 'vectors', 'assistants', 'geocoding']
    }
];

const docsLinks = [
    {label: 'Marketplace listing pages', href: '/marketplace'},
    {label: 'Social page', href: '/social'},
    {label: 'Sell page', href: '/sell'},
    {label: 'API docs', href: '/docs/api'},
    {label: 'Vendor guide', href: '/docs/vendor-guide'},
    {label: 'Enterprise brochure', href: '/resources/enterprise-brochure'},
    {label: 'Whitepaper', href: '/resources/whitepaper'}
];

function getPageCopy(locale) {
    return pageCopy[locale] || pageCopy.en;
}

export async function generateMetadata({params}) {
    const {locale} = await params;
    const copy = getPageCopy(locale);

    return {
        title: copy.title,
        description: copy.lead
    };
}

export default async function SfxMarketSubsitePage({params}) {
    const {locale} = await params;
    const copy = getPageCopy(locale);

    return (
        <div className={styles.page}>
            <div className={styles.container}>
                <header className={styles.panel}>
                    <span className={styles.badge}>{copy.badge}</span>
                    <h1 className={styles.title}>{copy.title}</h1>
                    <p className={styles.lead}>{copy.lead}</p>
                    <div className={styles.ctaRow}>
                        <Link className={styles.primaryButton} href="/marketplace">{copy.ctaMarketplace}</Link>
                        <Link className={styles.secondaryButton} href="/social">{copy.ctaSocial}</Link>
                        <Link className={styles.secondaryButton} href="/docs/api">{copy.ctaDocs}</Link>
                    </div>
                </header>

                <section className={styles.panel}>
                    <h2 className={styles.sectionTitle}>{copy.fitTitle}</h2>
                    <ul className={styles.bulletList}>
                        {copy.fitPoints.map((point) => (
                            <li key={point}>{point}</li>
                        ))}
                    </ul>
                </section>

                <section className={styles.panel}>
                    <h2 className={styles.sectionTitle}>{copy.examplesTitle}</h2>
                    <div className={styles.gridTwo}>
                        {copy.examples.map((item) => (
                            <article key={item.title} className={styles.card}>
                                <h3 className={styles.cardTitle}>{item.title}</h3>
                                <p className={styles.cardText}>{item.text}</p>
                            </article>
                        ))}
                    </div>
                </section>

                <section className={styles.panel}>
                    <h2 className={styles.sectionTitle}>{copy.servicesTitle}</h2>
                    <p className={styles.sectionLead}>{copy.servicesIntro}</p>
                    <div className={styles.gridTwo}>
                        {serviceGroups.map((group) => (
                            <article key={group.title} className={styles.card}>
                                <h3 className={styles.cardTitle}>{group.title}</h3>
                                <ul className={styles.moduleList}>
                                    {group.modules.map((moduleName) => (
                                        <li key={moduleName} className={styles.moduleItem}>{moduleName}</li>
                                    ))}
                                </ul>
                            </article>
                        ))}
                    </div>
                </section>

                <section className={styles.panel}>
                    <h2 className={styles.sectionTitle}>{copy.architectureTitle}</h2>
                    <div className={styles.gridTwo}>
                        {copy.architectureItems.map((item) => (
                            <article key={item.title} className={styles.card}>
                                <h3 className={styles.cardTitle}>{item.title}</h3>
                                <p className={styles.cardText}>{item.text}</p>
                            </article>
                        ))}
                    </div>
                </section>

                <section className={styles.panel}>
                    <h2 className={styles.sectionTitle}>{copy.valueTitle}</h2>
                    <ul className={styles.bulletList}>
                        {copy.valuePoints.map((point) => (
                            <li key={point}>{point}</li>
                        ))}
                    </ul>
                </section>

                <section className={styles.panel}>
                    <h2 className={styles.sectionTitle}>{copy.docsTitle}</h2>
                    <div className={styles.docsGrid}>
                        {docsLinks.map((item) => (
                            <Link key={item.href} href={item.href} className={styles.docLink}>
                                {item.label}
                            </Link>
                        ))}
                    </div>
                    <p className={styles.note}>{copy.finalNote}</p>
                </section>
            </div>
        </div>
    );
}
