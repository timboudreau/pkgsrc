# $NetBSD: Makefile,v 1.7 2019/03/16 23:56:02 wen Exp $

DISTNAME=	Sub-Quote-2.006003
PKGNAME=	p5-${DISTNAME}
CATEGORIES=	devel perl5
MASTER_SITES=	${MASTER_SITE_PERL_CPAN:=Sub/}

MAINTAINER=	pkgsrc-users@NetBSD.org
HOMEPAGE=	http://search.cpan.org/dist/Sub-Quote/
COMMENT=	Efficient generation of subroutines via string eval
LICENSE=	${PERL5_LICENSE}

BUILD_DEPENDS+=	p5-Test-Fatal>=0.003:../../devel/p5-Test-Fatal

PERL5_PACKLIST=	auto/Sub/Quote/.packlist

.include "../../lang/perl5/module.mk"
.include "../../mk/bsd.pkg.mk"
