# $NetBSD: Makefile,v 1.44 2019/04/03 00:32:31 ryoon Exp $

GITHUB_PROJECT=	sqlite3-ocaml
GITHUB_RELEASE=	${PKGVERSION_NOREV}
DISTNAME=	sqlite3-4.4.1
PKGNAME=	ocaml-${DISTNAME}
PKGREVISION=	3
CATEGORIES=	databases
MASTER_SITES=	${MASTER_SITE_GITHUB:=mmottl/}
EXTRACT_SUFX=	.tbz

MAINTAINER=	jaapb@NetBSD.org
HOMEPAGE=	https://github.com/mmottl/sqlite3-ocaml/
COMMENT=	Sqlite bindings for OCaml
LICENSE=	mit

USE_TOOLS+=	gmake pkg-config
OCAML_USE_DUNE=	yes

.include "../../mk/ocaml.mk"
.include "../../databases/sqlite3/buildlink3.mk"
.include "../../devel/ocaml-base/buildlink3.mk"
.include "../../devel/ocaml-configurator/buildlink3.mk"
.include "../../devel/ocaml-stdio/buildlink3.mk"
.include "../../mk/bsd.pkg.mk"
