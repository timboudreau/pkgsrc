# $NetBSD: Makefile,v 1.45 2019/04/20 12:32:26 wen Exp $

DISTNAME=		Net-DNS-SEC-1.12
PKGNAME=		p5-${DISTNAME}
CATEGORIES=		security net perl5
MASTER_SITES=		${MASTER_SITE_PERL_CPAN:=Net/}

MAINTAINER=		pkgsrc-users@NetBSD.org
HOMEPAGE=		http://www.net-dns.org/
COMMENT=		DNSSEC extension to Net::DNS Perl DNS Resolver Module
LICENSE=		modified-bsd

DEPENDS+=		p5-Digest-BubbleBabble>=0.01:../../security/p5-Digest-BubbleBabble
DEPENDS+=		p5-Digest-SHA1>=2.04:../../security/p5-Digest-SHA1
DEPENDS+=		p5-MIME-Base32-[0-9]*:../../converters/p5-MIME-Base32
DEPENDS+=		p5-Net-DNS>=0.64:../../net/p5-Net-DNS

#BUILDLINK_API_DEPENDS.perl+=	perl>=5.6.1

PERL5_PACKLIST=		auto/Net/DNS/SEC/.packlist

REPLACE_PERL=		demo/getkeyset demo/key2ds demo/make-signed-keyset

DOCSUBDIR=		share/doc/${PKGBASE}
EGSUBDIR=		share/examples/${PKGBASE}

INSTALLATION_DIRS+=	${DOCSUBDIR} ${EGSUBDIR}

SUBST_CLASSES+=		p5nds
SUBST_STAGE.p5nds=	post-install
SUBST_FILES.p5nds=	demo/getkeyset
SUBST_SED.p5nds=	-e "1s,-I./blib/lib,,"

post-install:
	${INSTALL_DATA} ${WRKSRC}/README  ${DESTDIR}${PREFIX}/${DOCSUBDIR}
.for f in getkeyset key2ds make-signed-keyset
	${INSTALL_SCRIPT} ${WRKSRC}/demo/${f} ${DESTDIR}${PREFIX}/${EGSUBDIR}
.endfor

.include "../../lang/perl5/module.mk"
.include "../../security/openssl/buildlink3.mk"
.include "../../mk/bsd.pkg.mk"
