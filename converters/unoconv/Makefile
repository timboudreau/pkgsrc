# $NetBSD: Makefile,v 1.3 2019/04/25 07:32:40 maya Exp $

DISTNAME=	unoconv-0.7
PKGREVISION=	1
CATEGORIES=	converters
MASTER_SITES=	${MASTER_SITE_GITHUB:=dagwieers/}

MAINTAINER=	pkgsrc-users@NetBSD.org
HOMEPAGE=	https://github.com/dagwieers/unoconv/
COMMENT=	Convert any format supported by LibreOffice/OpenOffice
LICENSE=	gnu-gpl-v2

USE_TOOLS+=	gmake

REPLACE_PYTHON+=unoconv

BUILD_DEPENDS+=	asciidoc-[0-9]*:../../textproc/asciidoc
DEPENDS+=	libreoffice>=5.3.4.2nb1:../../misc/libreoffice

pre-build:
	sed "s|@PREFIX@|${PREFIX}|" ${FILESDIR}/unoconv.sh > \
	    ${WRKDIR}/unoconv.sh

do-install:
	${INSTALL_PROGRAM_DIR} ${DESTDIR}/${PREFIX}/bin
	${INSTALL_PROGRAM_DIR} ${DESTDIR}/${PREFIX}/libexec
	${INSTALL_MAN_DIR} ${DESTDIR}/${PREFIX}/${PKGMANDIR}/man1
	${INSTALL_SCRIPT} ${WRKSRC}/unoconv \
	    ${DESTDIR}${PREFIX}/libexec/unoconv.py
	${INSTALL_SCRIPT} ${WRKDIR}/unoconv.sh \
	    ${DESTDIR}${PREFIX}/bin/unoconv
	${INSTALL_MAN} ${WRKSRC}/doc/unoconv.1 \
	    ${DESTDIR}${PREFIX}/${PKGMANDIR}/man1/unoconv.1

.include "../../lang/python/application.mk"
.include "../../mk/bsd.pkg.mk"
