# $NetBSD: Makefile,v 1.3 2019/04/08 15:48:31 adam Exp $

DISTNAME=	certbot-dns-nsone-0.33.1
PKGNAME=	${PYPKGPREFIX}-${DISTNAME}
CATEGORIES=	security python
MASTER_SITES=	${MASTER_SITE_PYPI:=c/certbot-dns-nsone/}

MAINTAINER=	triaxx@NetBSD.org
HOMEPAGE=	https://github.com/certbot/certbot
COMMENT=	NS1 DNS Authenticator plugin for Certbot
LICENSE=	apache-2.0

DEPENDS+=	${PYPKGPREFIX}-acme>=0.31.0:../../security/py-acme
DEPENDS+=	${PYPKGPREFIX}-certbot>=0.31.0:../../security/py-certbot
DEPENDS+=	${PYPKGPREFIX}-lexicon>=2.2.1:../../net/py-lexicon
DEPENDS+=	${PYPKGPREFIX}-mock-[0-9]*:../../devel/py-mock
DEPENDS+=	${PYPKGPREFIX}-setuptools-[0-9]*:../../devel/py-setuptools
DEPENDS+=	${PYPKGPREFIX}-ZopeInterface-[0-9]*:../../devel/py-ZopeInterface

USE_LANGUAGES=	# none

.include "../../lang/python/egg.mk"
.include "../../mk/bsd.pkg.mk"
